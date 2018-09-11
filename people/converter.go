package people

import (
	"strings"

	"github.com/Financial-Times/neo-model-utils-go/mapper"
)

const (
	thingsApiUrl = "http://api.ft.com/things/"
	ftThing      = "http://www.ft.com/thing/"
)

func convertToPerson(concept Concept, p *Person) {
	p.ID = convertID(concept.ID)
	p.APIURL = convertApiUrl(concept.APIURL, "people")
	p.PrefLabel = concept.PrefLabel
	p.Description = concept.Description
	p.DescriptionXML = concept.DescriptionXML
	p.ImageURL = concept.ImageURL
	p.Salutation = concept.Salutation
	p.BirthYear = concept.BirthYear
	p.Types = mapper.FullTypeHierarchy(concept.Type)
	p.DirectType = concept.Type

	for _, account := range concept.Account {
		switch {
		case strings.Contains(account.Type, "facebookProfile"):
			p.FacebookProfile = account.Value.(string)
		case strings.Contains(account.Type, "twitterHandle"):
			p.TwitterHandle = account.Value.(string)
		case strings.Contains(account.Type, "emailAddress"):
			p.EmailAddress = account.Value.(string)
		}
	}

	var labels []string
	for _, label := range concept.AlternativeLabels {
		labels = append(labels, label.Value.(string))
	}
	p.Labels = labels

	var memberships []Membership
	for _, related := range concept.RelatedConcepts {
		memberships = append(memberships, *convertToMembership(related.Concept))
	}
	p.Memberships = memberships
}

func convertToMembership(c Concept) *Membership {
	var organisation Organisation
	for _, related := range c.RelatedConcepts {
		if strings.Contains(related.Concept.Type, "Organisation") {
			organisation = *convertToOrganisation(related.Concept)
			break
		}
	}

	var roles []Role
	for _, role := range c.RelatedConcepts {
		if strings.Contains(role.Concept.Type, "Role") {
			roles = append(roles, *convertToRole(role.Concept))
		}
	}

	var m Membership
	m.Title = c.PrefLabel
	m.Types = mapper.FullTypeHierarchy(c.Type)
	m.DirectType = c.Type
	m.Organisation = organisation
	m.Roles = roles
	if len(c.ChangeEvents) > 0 {
		m.ChangeEvents = c.ChangeEvents
	}

	return &m
}

func convertToOrganisation(c Concept) *Organisation {
	var o Organisation
	o.ID = convertID(c.ID)
	o.APIURL = convertApiUrl(c.APIURL, "organisations")
	o.PrefLabel = c.PrefLabel
	o.Types = mapper.FullTypeHierarchy(c.Type)
	o.DirectType = c.Type
	return &o
}

func convertToRole(c Concept) *Role {
	var r Role
	r.ID = convertID(c.ID)
	r.APIURL = convertApiUrl(c.APIURL, "things")
	r.PrefLabel = c.PrefLabel
	r.Types = mapper.FullTypeHierarchy(c.Type)
	r.DirectType = c.Type

	if len(c.ChangeEvents) > 0 {
		r.ChangeEvents = c.ChangeEvents
	}

	return &r
}

func convertApiUrl(conceptsApiUrl string, desired string) string {
	return strings.Replace(conceptsApiUrl, "concepts", desired, 1)
}

func convertID(conceptsApiID string) string {
	return strings.Replace(conceptsApiID, ftThing, thingsApiUrl, 1)
}
