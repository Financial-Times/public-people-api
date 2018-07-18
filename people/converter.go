package people

import (
	"github.com/Financial-Times/neo-model-utils-go/mapper"
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
	p.Types = mapper.FullTypeHierarchy(concept.Type)
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

	var m Membership
	m.Title = c.PrefLabel
	m.Types = mapper.FullTypeHierarchy(c.Type)
	m.DirectType = c.Type
	m.Organisation = organisation
	m.Roles = roles
	changeEvents := getChangeEvengts(c)
	if len(changeEvents) > 0 {
		m.ChangeEvents = changeEvents
	}

	return &m
}

func convertToOrganization(c Concept) *Organisation {
	var o Organisation
	o.ID = c.ID
	o.APIURL = c.APIURL
	o.PrefLabel = c.PrefLabel
	o.Types = mapper.FullTypeHierarchy(c.Type)
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
	r.Types = mapper.FullTypeHierarchy(c.Type)
	r.DirectType = c.Type

	changeEvents := getChangeEvengts(c)
	if len(changeEvents) > 0 {
		r.ChangeEvents = changeEvents
	}

	return &r
}

func getChangeEvengts(c Concept) []ChangeEvent {
	if len(c.ChangeEvents) > 0 {
		return c.ChangeEvents
	}

	var changeEvents []ChangeEvent
	if len(c.InceptionDate) > 0 {
		changeEvents = append(changeEvents, ChangeEvent{
			StartedAt: c.InceptionDate,
		})
	}
	if len(c.TerminationDate) > 0 {
		changeEvents = append(changeEvents, ChangeEvent{
			EndedAt: c.TerminationDate,
		})
	}
	return changeEvents
}
