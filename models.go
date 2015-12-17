package main

import (
	"fmt"
	"log"
	"time"
)

type neoPerson struct {
	UUID      string `json:"uuid"`
	PrefLabel string `json:"prefLabel"`
	Name      string `json:"Name"`
}

// Common structure TODO get Embedding working
type Common struct {
	ID        string `json:"id"`
	APIURL    string `json:"apiUrl"`
	PrefLabel string `json:"prefLabel"`
}

// Person structure for writing to responses
type Person struct {
	ID          string       `json:"id"`
	APIURL      string       `json:"apiUrl"`
	PrefLabel   string       `json:"prefLabel"`
	Types       []string     `json:"types"`
	Labels      []string     `json:"labels"`
	Profile     string       `json:"profile"`
	Memberships []Membership `json:"memberships"`
}

// TODO maybe return err too ?
func toPerson(neoPerson neoPerson) (person Person) {
	log.Printf("Incoming neoPerson %v+", neoPerson)
	person.APIURL = fmt.Sprintf("http://api.ft.com/people/%s", neoPerson.UUID)
	person.ID = fmt.Sprintf("http://api.ft.com/things/%s", neoPerson.UUID)
	if neoPerson.PrefLabel != "" {
		person.PrefLabel = neoPerson.PrefLabel
	} else {
		person.PrefLabel = neoPerson.Name
	}
	person.Types = []string{"Person"}
	log.Printf("Outgoing Person %v+", person)
	return person
}

// Membership structure
type Membership struct {
	Title        string        `json:"title"`
	Organisation Organisation  `json:"organisation"`
	Roles        []Role        `json:"roles"`
	ChangeEvents []ChangeEvent `json:"changeEvents"`
}

// Organisation structure
type Organisation struct {
	ID        string   `json:"id"`
	APIURL    string   `json:"apiUrl"`
	PrefLabel string   `json:"prefLabel"`
	Types     []string `json:"types"`
}

// Role structure
type Role struct {
	ID           string        `json:"id"`
	APIURL       string        `json:"apiUrl"`
	PrefLabel    string        `json:"prefLabel"`
	ChangeEvents []ChangeEvent `json:"changeEvents"`
}

// ChangeEvent structure TODO prevent 'zero' values being encoded
// http://stackoverflow.com/questions/18088294/how-to-not-marshal-an-empty-struct-into-json-with-go
type ChangeEvent struct {
	StartedAt *time.Time `json:"startedAt,omitempty"`
	EndedAt   *time.Time `json:"endedAt,omitempty"`
}
