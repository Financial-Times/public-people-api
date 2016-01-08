package main

// Thing is the base entity, all nodes in neo4j should have these properties
/* The following is currently defined in Java (3da1b900b38)
@JsonInclude(NON_EMPTY)
public class Thing {
    public String id;
    public String apiUrl;
    public String prefLabel;
    public List<String> types = new ArrayList<>();
}
*/
type Thing struct {
	ID        string   `json:"id"`
	APIURL    string   `json:"apiUrl"` // self ?
	PrefLabel string   `json:"prefLabel,omitempty"`
	Types     []string `json:"types"`
}

// Person is the structure used for the people API
/* The following is currently defined in Java (e4b93668e32)
@JsonInclude(NON_EMPTY)
public class Person extends Thing {
    public List<String> labels = new ArrayList<>();
    public String profile;
    public List<Membership> memberships = new ArrayList<>();
}
*/
type Person struct {
	*Thing
	Labels      *[]string    `json:"labels,omitempty"`
	Memberships []Membership `json:"memberships"`
	Salutation  string       `json:"salutation,omitempty"`
	BirthYear   string       `json:"birthYear,omitempty"`
}

// Membership represents the relationship between a person and their roles associated with an organisation
/*
@JsonInclude(Include.NON_EMPTY)
public class Membership {
    public String title;
    public Thing organisation;
    public Thing person;
    public List<ChangeEvent> changeEvents = new ArrayList();
    public List<MembershipRole> roles = new ArrayList();
*/
type Membership struct {
	Title        string       `json:"title,omitempty"`
	PrefLabel    string       `json:"title,omitempty"`
	Organisation Organisation `json:"organisation"`
	Roles        []Role       `json:"roles"`
	ChangeEvent  *ChangeEvent `json:"changeEvent,omitempty"`
}

// Organisation simplified representation used in Person API
type Organisation struct {
	*Thing
	Labels *[]string `json:"labels,omitempty"`
}

// Role represents the capacity or funciton that a person performs for an organisation
/*
@JsonInclude(Include.NON_EMPTY)
public class MembershipRole extends Thing {
    public List<ChangeEvent> changeEvents = new ArrayList();
}
*/
type Role struct {
	*Thing
	ChangeEvent *ChangeEvent `json:"changeEvent,omitempty"`
}

// ChangeEvent represents when something started or ended
/*
@JsonInclude(Include.NON_EMPTY)
public class ChangeEvent {
    public String startedAt;
    public String endedAt;
*/
type ChangeEvent struct {
	Started string `json:"started,omitempty"`
	Ended   string `json:"ended,omitempty"`
}
