package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jmcvetta/neoism"
)

// PeopleDriver interface
type PeopleDriver interface {
	Read(id string) (person Person, found bool, err error)
}

// PeopleCypherDriver struct
type PeopleCypherDriver struct {
	db *neoism.Database
}

//NewPeopleCypherDriver instanciate driver
func NewPeopleCypherDriver(db *neoism.Database) PeopleCypherDriver {
	return PeopleCypherDriver{db}
}

func (pcw PeopleCypherDriver) Read(uuid string) (person Person, found bool, err error) {
	results := []struct {
		Rs []struct {
			P         neoism.Node
			PrefLabel string
			// P struct {
			// 	PrefLabel string `json="data.prefLabel"`
			// }
		}
	}{}
	query := &neoism.CypherQuery{
		Statement: `
                MATCH (p:Person{uuid: {uuid} })<-[:HAS_MEMBER]-(m:Membership)
                OPTIONAL MATCH (m)-[:HAS_ORGANISATION]->(o:Organisation)
                OPTIONAL MATCH (m)-[rr:HAS_ROLE]->(r:Role)
                WITH p, collect({m:m, o:o, r:r, rr:rr}) as m
                RETURN collect ({p:p, prefLabel:p.prefLabel, memberships:m}) as rs
                        `,
		Parameters: neoism.Props{"uuid": uuid},
		Result:     &results,
	}
	err = pcw.db.Cypher(query)
	if err != nil {
		log.Errorf("Error from neoism %+v", err)
		return Person{}, false, err
	}
	fmt.Printf("Neoism %+v\n", results[0].Rs[0])
	//results[0].P.Db = pcw.db
	//Thing(results[0].P, &result)
	log.Debugf("Returning %v", person)
	return person, true, nil
}
