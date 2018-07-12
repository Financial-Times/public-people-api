package people

import (
	"strings"
)

func convertToPerson(concept Concept, p *Person) {
	p.ID = concept.ID
	p.APIURL = concept.APIURL
	p.PrefLabel = concept.PrefLabel
	p.Description = concept.Description
	p.DescriptionXML = concept.descriptionXML
	p.ImageURL = concept.ImageURL
	p.Salutation = concept.Salutation
	p.BirthYear = concept.BirthYear
	p.Types = []string{
		"http://www.ft.com/ontology/core/Thing",
		"http://www.ft.com/ontology/concept/Concept",
		"http://www.ft.com/ontology/person/Person",
	}
	p.DirectType = concept.Type

	for _, account := range concept.Account {
		switch {
		case strings.
			Contains(account.Type, "facebookProfile"):
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
			organisation = *convertToOrganization(related.Concept)
			break
		}
	}

	var roles []Role
	for _, role := range c.RelatedConcepts {
		if strings.Contains(role.Concept.Type, "Role") {
			roles = append(roles, *convertToRole(role.Concept))
		}
	}

	return &Membership{
		Title: c.PrefLabel,
		Types: []string{
			"http://www.ft.com/ontology/core/Thing",
			"http://www.ft.com/ontology/concept/Concept",
			"http://www.ft.com/ontology/organisation/Membership",
		},
		DirectType:   c.Type,
		Organisation: organisation,
		ChangeEvents: getChangeEvents(c),
		Roles:        roles,
	}
}

func convertToOrganization(c Concept) *Organisation {
	var o Organisation
	o.ID = c.ID
	o.APIURL = c.APIURL
	o.PrefLabel = c.PrefLabel
	o.Types = []string{
		"http://www.ft.com/ontology/core/Thing",
		"http://www.ft.com/ontology/concept/Concept",
		"http://www.ft.com/ontology/organisation/Organisation",
	}
	o.DirectType = c.Type

	var labels []string
	for _, label := range c.AlternativeLabels {
		labels = append(labels, label.Value.(string))
	}
	o.Labels = labels

	return &o
}

func convertToRole(c Concept) *Role {
	var r Role
	r.ID = c.ID
	r.APIURL = c.APIURL
	r.PrefLabel = c.PrefLabel
	r.Types = []string{
		"http://www.ft.com/ontology/core/Thing",
		"http://www.ft.com/ontology/concept/Concept",
		"http://www.ft.com/ontology/MembershipRole",
	}
	r.DirectType = c.Type
	r.ChangeEvents = getChangeEvents(c)

	return &r
}

func getChangeEvents(c Concept) []ChangeEvent {
	return []ChangeEvent{
		ChangeEvent{
			StartedAt: c.InceptionDate,
			EndedAt:   c.TerminationDate,
		},
	}
}
