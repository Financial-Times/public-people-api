package people

import (
	"errors"
	"fmt"
	"time"

	"github.com/Financial-Times/neo-model-utils-go/mapper"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	log "github.com/Sirupsen/logrus"
	"github.com/jmcvetta/neoism"
	"github.com/satori/go.uuid"
)

// Driver interface
type Driver interface {
	Read(id uuid.UUID) (person Person, found bool, err error)
	CheckConnectivity() error
}

// CypherDriver struct
type CypherDriver struct {
	conn neoutils.NeoConnection
	env  string
}

//NewCypherDriver instantiate driver
func NewCypherDriver(conn neoutils.NeoConnection, env string) CypherDriver {
	return CypherDriver{conn, env}
}

// CheckConnectivity tests neo4j by running a simple cypher query
func (pcw CypherDriver) CheckConnectivity() error {
	return neoutils.Check(pcw.conn)
}

type neoChangeEvent struct {
	StartedAt string
	EndedAt   string
}

type neoReadStruct struct {
	P struct {
		ID             string
		Types          []string
		PrefLabel      string
		Labels         []string
		Salutation     string
		BirthYear      int
		EmailAddress   string
		TwitterHandle  string
		Description    string
		DescriptionXML string
		ImageURL       string
	}
	M []struct {
		M struct {
			ID           string
			Types        []string
			PrefLabel    string
			Title        string
			ChangeEvents []neoChangeEvent
		}
		O struct {
			ID        string
			Types     []string
			PrefLabel string
			Labels    []string
		}
		R []struct {
			ID           string
			Types        []string
			PrefLabel    string
			ChangeEvents []neoChangeEvent
		}
	}
}

func (pcw CypherDriver) Read(uuid uuid.UUID) (person Person, found bool, err error) {
	person = Person{}
	results := []struct {
		Rs []neoReadStruct
	}{}
	sixMonthsEpoch := time.Now().Unix() - (60 * 60 * 24 * 30 * 6)
	query := &neoism.CypherQuery{
		Statement: `
                        MATCH (identifier:UPPIdentifier{value:{uuid}})
                        MATCH (identifier)-[:IDENTIFIES]->(p:Person)
                        OPTIONAL MATCH (p)<-[:HAS_MEMBER]-(m:Membership)
                        OPTIONAL MATCH (m)-[:HAS_ORGANISATION]->(o:Organisation)
                        OPTIONAL MATCH (m)-[rr:HAS_ROLE]->(r:Role)
                        OPTIONAL MATCH (o)<-[rel:MENTIONS]-(c:Content) WHERE c.publishedDateEpoch > {publishedDateEpoch}
                        WITH    p,
                                { id:o.uuid, types:labels(o), prefLabel:o.prefLabel, annCount:COUNT(c)} as o,
                                { id:m.uuid, types:labels(m), prefLabel:m.prefLabel, title:m.title, changeEvents:[{startedAt:m.inceptionDate}, {endedAt:m.terminationDate}] } as m,
                                { id:r.uuid, types:labels(r), prefLabel:r.prefLabel, changeEvents:[{startedAt:rr.inceptionDate}, {endedAt:rr.terminationDate}] } as r
                        WITH p, m, o, collect(r) as r ORDER BY o.annCount DESC
                        WITH p, collect({m:m, o:o, r:r}) as m
                        WITH m, { id:p.uuid, types:labels(p), prefLabel:p.prefLabel, labels:p.aliases,
												     birthYear:p.birthYear, salutation:p.salutation, emailAddress:p.emailAddress,
														 twitterHandle:p.twitterHandle, imageURL:p.imageURL,
														 Description:p.description, descriptionXML:p.descriptionXML} as p
                        RETURN collect ({p:p, m:m}) as rs
                        `,
		Parameters: neoism.Props{"uuid": uuid.String(), "publishedDateEpoch": sixMonthsEpoch},
		Result:     &results,
	}

	if err := pcw.conn.CypherBatch([]*neoism.CypherQuery{query}); err != nil || len(results) == 0 || len(results[0].Rs) == 0 {
		return Person{}, false, err
	} else if len(results) != 1 && len(results[0].Rs) != 1 {
		errMsg := fmt.Sprintf("Multiple people found with the same uuid:%s !", uuid)
		log.Error(errMsg)
		return Person{}, true, errors.New(errMsg)
	}

	person = neoReadStructToPerson(results[0].Rs[0], pcw.env)
	return person, true, nil
}

func neoReadStructToPerson(neo neoReadStruct, env string) Person {
	public := Person{}
	public.Thing = &Thing{}
	public.ID = mapper.IDURL(neo.P.ID)
	public.APIURL = mapper.APIURL(neo.P.ID, neo.P.Types, env)
	public.Types = mapper.TypeURIs(neo.P.Types)
	public.PrefLabel = neo.P.PrefLabel
	if len(neo.P.Labels) > 0 {
		public.Labels = &neo.P.Labels
	}
	public.BirthYear = neo.P.BirthYear
	public.Salutation = neo.P.Salutation
	public.Description = neo.P.Description
	public.DescriptionXML = neo.P.DescriptionXML
	public.EmailAddress = neo.P.EmailAddress
	public.TwitterHandle = neo.P.TwitterHandle
	public.ImageURL = neo.P.ImageURL

	if len(neo.M) == 1 && (neo.M[0].M.ID == "") {
		public.Memberships = make([]Membership, 0, 0)
	} else {
		public.Memberships = make([]Membership, len(neo.M))
		for mIdx, neoMem := range neo.M {
			membership := Membership{}
			membership.Title = neoMem.M.PrefLabel
			membership.Organisation = Organisation{}
			membership.Organisation.Thing = &Thing{}
			membership.Organisation.ID = mapper.IDURL(neoMem.O.ID)
			membership.Organisation.APIURL = mapper.APIURL(neoMem.O.ID, neoMem.O.Types, env)
			membership.Organisation.Types = mapper.TypeURIs(neoMem.O.Types)
			membership.Organisation.PrefLabel = neoMem.O.PrefLabel
			if len(neoMem.O.Labels) > 0 {
				membership.Organisation.Labels = &neoMem.O.Labels
			}
			if a, b := changeEvent(neoMem.M.ChangeEvents); a == true {
				membership.ChangeEvents = b
			}
			membership.Roles = make([]Role, len(neoMem.R))
			for rIdx, neoRole := range neoMem.R {
				role := Role{}
				role.Thing = &Thing{}
				role.ID = mapper.IDURL(neoRole.ID)
				role.APIURL = mapper.APIURL(neoRole.ID, neoRole.Types, env)
				role.PrefLabel = neoRole.PrefLabel
				if a, b := changeEvent(neoRole.ChangeEvents); a == true {
					role.ChangeEvents = b
				}

				membership.Roles[rIdx] = role
			}
			public.Memberships[mIdx] = membership
		}
	}
	return public
}

func changeEvent(neoChgEvts []neoChangeEvent) (bool, *[]ChangeEvent) {
	var results []ChangeEvent
	currentLayout := "2006-01-02T15:04:05.999Z"
	layout := "2006-01-02T15:04:05Z"

	if neoChgEvts[0].StartedAt == "" && neoChgEvts[1].EndedAt == "" {
		results = make([]ChangeEvent, 0, 0)
		return false, &results
	}
	for _, neoChgEvt := range neoChgEvts {
		if neoChgEvt.StartedAt != "" {
			t, _ := time.Parse(currentLayout, neoChgEvt.StartedAt)
			results = append(results, ChangeEvent{StartedAt: t.Format(layout)})
		}
		if neoChgEvt.EndedAt != "" {
			t, _ := time.Parse(layout, neoChgEvt.EndedAt)
			results = append(results, ChangeEvent{EndedAt: t.Format(layout)})
		}
	}
	return true, &results
}
