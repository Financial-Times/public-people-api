package main

// Thing is the base entity, all nodes in neo4j should have these properties
type Thing struct {
	ID     string   `json:"id"`
	APIURL string   `json:"apiUrl"` // self ?
	Types  []string `json:"type"`
}

// Person is the structure used for the people API
type Person struct {
	ID          string       `json:"id"`
	APIURL      string       `json:"apiUrl"` // self ?
	Types       []string     `json:"type"`
	PrefLabel   string       `json:"prefLabel"`
	Name        string       `json:"name"`
	Salutation  string       `json:"salutation"`
	BirthYear   string       `json:"birthYear"`
	Memberships []Membership `json:"memberships"`
}

// Membership represents the relationship between a person and their roles assoicated with an organisation
type Membership struct {
	ID           string                 `json:"id"`
	APIURL       string                 `json:"apiUrl"` // self ?
	Types        []string               `json:"type"`
	Organisation MembershipOrganisation `json:"organisation"`
	Roles        []Role                 `json:"roles"`
	ChangeEvent  ChangeEvent            `json:"changeEvent"`
}

// MembershipOrganisation represents an Organisation of some sort, for example a company or educational establishment
type MembershipOrganisation struct {
	ID     string   `json:"id"`
	APIURL string   `json:"apiUrl"` // self ?
	Types  []string `json:"type"`
}

// Role represents the capacity or funciton that a person performs for an organisation
type Role struct {
	ID          string      `json:"id"`
	APIURL      string      `json:"apiUrl"` // self ?
	Types       []string    `json:"type"`
	ChangeEvent ChangeEvent `json:"changeEvent"`
}

// ChangeEvent represents when something started or ended
type ChangeEvent struct {
	Started string `json:"started,omitempty"`
	Ended   string `json:"ended,omitempty"`
}
