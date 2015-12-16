package main

import (
	"time"
)

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
	StartedAt time.Time `json:"startedAt,omitempty"`
	EndedAt   time.Time `json:"endedAt,omitempty"`
}
