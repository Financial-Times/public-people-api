package people

// Thing is the base entity, all nodes in neo4j should have these properties
type Thing struct {
	ID        string `json:"id"`
	APIURL    string `json:"apiUrl"` // self ?
	PrefLabel string `json:"prefLabel,omitempty"`
}

// Person is the structure used for the people API
type Person struct {
	Thing
	Types           []string     `json:"types"`
	DirectType      string       `json:"directType,omitempty"`
	Labels          []string     `json:"labels,omitempty"`
	Memberships     []Membership `json:"memberships,omitempty"`
	Salutation      string       `json:"salutation,omitempty"`
	BirthYear       int          `json:"birthYear,omitempty"`
	EmailAddress    string       `json:"emailAddress,omitempty"`
	TwitterHandle   string       `json:"twitterHandle,omitempty"`
	FacebookProfile string       `json:"facebookProfile,omitempty"`
	Description     string       `json:"description,omitempty"`
	DescriptionXML  string       `json:"descriptionXML,omitempty"`
	ImageURL        string       `json:"_imageUrl,omitempty"` // TODO we should implement this properly as an imageset
	IsDeprecated    bool         `json:"isDeprecated,omitempty"`
}

// Membership represents the relationship between a person and their roles associated with an organisation
type Membership struct {
	Title        string        `json:"title,omitempty"`
	Types        []string      `json:"types"`
	DirectType   string        `json:"directType,omitempty"`
	Organisation Organisation  `json:"organisation"`
	ChangeEvents []ChangeEvent `json:"changeEvents,omitempty"`
	Roles        []Role        `json:"roles"`
}

// Organisation simplified representation used in Person API
type Organisation struct {
	Thing
	Types      []string `json:"types"`
	DirectType string   `json:"directType,omitempty"`
	Labels     []string `json:"labels,omitempty"`
}

// Role represents the capacity or funciton that a person performs for an organisation
type Role struct {
	Thing
	Types        []string      `json:"types"`
	DirectType   string        `json:"directType,omitempty"`
	ChangeEvents []ChangeEvent `json:"changeEvents,omitempty"`
}

// ChangeEvent represent when something started or ended
type ChangeEvent struct {
	StartedAt string `json:"startedAt,omitempty"`
	EndedAt   string `json:"endedAt,omitempty"`
}

type TypedValue struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type PredicateConcept struct {
	Concept   Concept `json:"concept"`
	Predicate string  `json:"predicate"`
}

type Concept struct {
	ID        string `json:"id"`
	APIURL    string `json:"apiUrl"`
	Type      string `json:"type"`
	PrefLabel string `json:"prefLabel"`
	// Additional fields
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"imageURL,omitempty"`
	// Account contains the values of:
	// * emailAddress
	// * facebookPage
	// * twitterHandle
	Account []TypedValue `json:"account,omitempty"`
	// AlternativeLabels contains the values of:
	// * aliases
	// * formerNames
	// * hiddenLabel
	// * legalName
	// * properName
	// * shortName
	// * tradeNames
	AlternativeLabels []TypedValue  `json:"alternativeLabels,omitempty"`
	ChangeEvents      []ChangeEvent `json:"changeEvents,omitempty"`
	// Person
	Salutation     string `json:"salutation,omitempty"`
	BirthYear      int    `json:"birthYear,omitempty"`
	DescriptionXML string `json:"descriptionXML,omitempty"`
	// Membership
	InceptionDate   string `json:"inceptionDate,omitempty"`
	TerminationDate string `json:"terminationDate,omitempty"`
	// Organisation
	CountryCode            string `json:"countryCode,omitempty"`
	CountryOfIncorporation string `json:"countryOfIncorporation,omitempty"`
	LeiCode                string `json:"leiCode,omitempty"`
	PostalCode             string `json:"postalCode,omitempty"`
	YearFounded            int    `json:"yearFounded,omitempty"`
	// Relations
	BroaderConcepts  []PredicateConcept `json:"broaderConcepts,omitempty"`
	NarrowerConcepts []PredicateConcept `json:"narrowerConcepts,omitempty"`
	RelatedConcepts  []PredicateConcept `json:"relatedConcepts,omitempty"`
	IsDeprecated     bool               `json:"isDeprecated,omitempty"`
}
