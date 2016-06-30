package people

// Thing is the base entity, all nodes in neo4j should have these properties
type Thing struct {
	ID        string `json:"id"`
	APIURL    string `json:"apiUrl"` // self ?
	PrefLabel string `json:"prefLabel,omitempty"`
}

// Person is the structure used for the people API
type Person struct {
	*Thing
	Types       []string     `json:"types"`
	Labels      *[]string    `json:"labels,omitempty"`
	Memberships []Membership `json:"memberships,omitempty"`
	Salutation  string       `json:"salutation,omitempty"`
	BirthYear   string       `json:"birthYear,omitempty"`
}

// Membership represents the relationship between a person and their roles associated with an organisation
type Membership struct {
	Title        string         `json:"title,omitempty"`
	Organisation Organisation   `json:"organisation"`
	ChangeEvents *[]ChangeEvent `json:"changeEvents,omitempty"`
	Roles        []Role         `json:"roles"`
}

// Organisation simplified representation used in Person API
type Organisation struct {
	*Thing
	Types  []string  `json:"types"`
	Labels *[]string `json:"labels,omitempty"`
}

// Role represents the capacity or funciton that a person performs for an organisation
type Role struct {
	*Thing
	ChangeEvents *[]ChangeEvent `json:"changeEvents,omitempty"`
}

// ChangeEvent represent when something started or ended
type ChangeEvent struct {
	StartedAt string `json:"startedAt,omitempty"`
	EndedAt   string `json:"endedAt,omitempty"`
}
