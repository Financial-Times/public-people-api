package people

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"reflect"
	"sort"

	"time"

	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	"github.com/Financial-Times/concepts-rw-neo4j/concepts"
	"github.com/Financial-Times/memberships-rw-neo4j/memberships"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/Financial-Times/organisations-rw-neo4j/organisations"
	person "github.com/Financial-Times/people-rw-neo4j/people"
	"github.com/Financial-Times/roles-rw-neo4j/roles"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
)

const (
	personId          = "13a9d251-71db-467a-af2f-7e56a61c910a"
	testTransactionID = "test_tid"
)

//Reusable Neo4J connection
var db neoutils.NeoConnection

var publicPeopleDriver CypherDriver

//Concept Services
var conceptsDriver concepts.ConceptService

// Old Model Services:
var peopleDriver baseftrwapp.Service
var organisationsDriver baseftrwapp.Service
var membershipsDriver baseftrwapp.Service
var rolesDriver baseftrwapp.Service

func init() {
	// We are initialising a lot of constraints on an empty database therefore we need the database to be fit before
	// we run tests so initialising the service will create the constraints first
	conf := neoutils.DefaultConnectionConfig()
	conf.Transactional = false

	db, _ = neoutils.Connect(neoUrl(), conf)
	if db == nil {
		panic("Cannot connect to Neo4J")
	}

	publicPeopleDriver = NewCypherDriver(db, "prod")

	peopleDriver = person.NewCypherPeopleService(db)
	peopleDriver.Initialise()

	organisationsDriver = organisations.NewCypherOrganisationService(db)
	organisationsDriver.Initialise()

	membershipsDriver = memberships.NewCypherMembershipService(db)
	membershipsDriver.Initialise()

	rolesDriver = roles.NewCypherDriver(db)
	rolesDriver.Initialise()

	conceptsDriver = concepts.NewConceptService(db)
	conceptsDriver.Initialise()

	duration := 5 * time.Second
	time.Sleep(duration)
}

func neoUrl() string {
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}
	return url
}

// TestNeoReadStructToPersonMandatoryFields checks that mandatory fields are set even if they are empty or nil / null
func TestNeoReadStructToPersonMandatoryFields(t *testing.T) {
	expected := `{"id":"http://api.ft.com/things/","apiUrl":"http://api.ft.com/things/","types":null}`
	person := neoReadStructToPerson(neoReadStruct{}, "prod")
	personJSON, err := json.Marshal(person)
	assert := assert.New(t)
	assert.NoError(err, "Unable to marshal Person to JSON")
	assert.Equal(expected, string(personJSON))
}

func TestNeoReadStructToPersonIncludingMultipleMemberships(t *testing.T) {
	defer cleanDB(db, t, "e903861d-7709-4ab3-aeb4-4d272ac4d105", "d137a439-3efd-4820-9cab-c200031e3dd9", "8865b295-c1f1-442e-8972-eb100dc50292",
		"2802a267-aa96-4f68-897c-66e90d7d57e8", "ac4be3c3-6dc1-4966-9cc5-ac824780f631", "638fc0c1-c4d9-4be4-b6d9-c97a057e7d1b", "0ee8e7b7-bac9-4db1-b94b-5605ce1d2907",
		"195df1b2-7e04-4c70-a865-4361c71e9a6b", "7b00924d-6115-4126-9bb5-5e3cfdfc8114", personId)

	writeJSONToService(t, peopleDriver, fmt.Sprintf("./fixtures/oldModel/Person-Siobhan_Morden-%s.json", personId))
	writeJSONToService(t, organisationsDriver, "./fixtures/oldModel/Organisation-Parent_A-638fc0c1-c4d9-4be4-b6d9-c97a057e7d1b.json")
	writeJSONToService(t, organisationsDriver, "./fixtures/oldModel/Organisation-Child_Of_A-ac4be3c3-6dc1-4966-9cc5-ac824780f631.json")
	writeJSONToService(t, organisationsDriver, "./fixtures/oldModel/Organisation-Other-2802a267-aa96-4f68-897c-66e90d7d57e8.json")
	writeJSONToService(t, membershipsDriver, "./fixtures/oldModel/Membership-Siobhan_Morden-8865b295-c1f1-442e-8972-eb100dc50292.json")
	writeJSONToService(t, membershipsDriver, "./fixtures/oldModel/Membership-Siobhan_Morden-d137a439-3efd-4820-9cab-c200031e3dd9.json")
	writeJSONToService(t, membershipsDriver, "./fixtures/oldModel/Membership-Siobhan_Morden-e903861d-7709-4ab3-aeb4-4d272ac4d105.json")
	writeJSONToService(t, rolesDriver, "./fixtures/oldModel/Role-0ee8e7b7-bac9-4db1-b94b-5605ce1d2907.json")

	person := readJSONtoPerson(t, "./fixtures/outputJSON/MultipleMembershipsOldModelTest-Output.json")
	readConceptAndCompare(t, person, personId)
}

func TestNeoReadPersonWithCanonicalUPPID(t *testing.T) {
	defer cleanDB(db, t, personId)
	writeJSONToService(t, peopleDriver, fmt.Sprintf("./fixtures/oldModel/Person-Siobhan_Morden-%s.json", personId))

	person := Person{
		Thing: Thing{
			ID:        "http://api.ft.com/things/13a9d251-71db-467a-af2f-7e56a61c910a",
			APIURL:    "http://api.ft.com/people/13a9d251-71db-467a-af2f-7e56a61c910a",
			PrefLabel: "Siobhan Morden",
		},
		Types:           []string{"http://www.ft.com/ontology/concept/Concept", "http://www.ft.com/ontology/core/Thing", "http://www.ft.com/ontology/person/Person"},
		DirectType:      "http://www.ft.com/ontology/person/Person",
		Labels:          []string{"Siobhan J Morden", "Siobhan Morden"},
		Salutation:      "Ms.",
		BirthYear:       1974,
		EmailAddress:    "test@example.com",
		TwitterHandle:   "@something",
		FacebookProfile: "the-facebook-profile",
		Description:     "Some text",
		DescriptionXML:  "Some text containing <strong>markup</strong>",
		Memberships:     []Membership{},
	}

	readConceptAndCompare(t, person, personId)
}

func TestNeoReadPersonWithAlternateUPPID(t *testing.T) {
	alternativePersonId := "d755c384-c302-485c-b12e-ea3c6751a6b6"

	cleanDB(db, t, alternativePersonId, personId)

	writeJSONToService(t, peopleDriver, fmt.Sprintf("./fixtures/oldModel/Person-Siobhan_Morden-%s.json", personId))

	publicPeopleDriver := NewCypherDriver(db, "prod")
	person, found, err := publicPeopleDriver.Read(alternativePersonId, testTransactionID)
	person1, found, err := publicPeopleDriver.Read(personId, testTransactionID)

	assert.NoError(t, err, "Error on reading Person with alternative id: %v", alternativePersonId)
	assert.True(t, found, "Person not found in database")
	assert.NoError(t, err, "Error on reading Person: %v", personId)
	assert.Equal(t, person, person1)

	assert.True(t, reflect.DeepEqual(person, person1), "Retrieving person with UUID: %v differed from retieveing the same person via another alternative UUID: %v", personId, alternativePersonId)
}

// New Model and backwards compatibility tests
func TestNewModelWithFullyNewModelMembershipRelatedConcepts(t *testing.T) {
	//New model to new model fully Type Org and Person and Role

	defer cleanDB(db, t, "7d0738b1-0ea2-47cb-bb82-e86744b389f0", "184cbe9b-b630-40d5-a5d0-99ecabd7fd86", "7ceeafe5-9f9a-4315-b3da-a5b4b69c013a", "8cdff2ba-3062-471e-b98a-7ee961239cd2")
	writeJSONToConceptsService(t, "./fixtures/newModel/MembershipRole-SmartyPants-7d0738b1-0ea2-47cb-bb82-e86744b389f0.json")
	writeJSONToConceptsService(t, "./fixtures/newModel/Organisation-RooneyRoosters-184cbe9b-b630-40d5-a5d0-99ecabd7fd86.json")
	writeJSONToConceptsService(t, "./fixtures/newModel/Person-Shirley-Rooney-7ceeafe5-9f9a-4315-b3da-a5b4b69c013a.json")
	writeJSONToConceptsService(t, "./fixtures/newModel/Membership-8cdff2ba-3062-471e-b98a-7ee961239cd2.json")

	person := readJSONtoPerson(t, "./fixtures/outputJSON/TestNewModelWithFullyNewModelMembershipRelatedConcepts.json")
	readConceptAndCompare(t, person, "7ceeafe5-9f9a-4315-b3da-a5b4b69c013a")
}

// TODO: When we concord to Factset we will need to handle a mixture of old model and new model

func TestNewModelWithThingOnlyMembershipRelatedConceptsDoesNotReturnMembership(t *testing.T) {
	// New model to org/person/role that is only a Thing - No membership should be returned

	defer cleanDB(db, t, "ef0921e4-c862-43ac-8936-f345b9fb131a", "7ceeafe5-9f9a-4315-b3da-a5b4b69c013a", "0ee8e7b7-bac9-4db1-b94b-5605ce1d2907", "ac4be3c3-6dc1-4966-9cc5-ac824780f631")

	writeJSONToService(t, membershipsDriver, "./fixtures/oldModel/Membership-Shirley-Rooney-ef0921e4-c862-43ac-8936-f345b9fb131a.json")
	writeJSONToConceptsService(t, "./fixtures/newModel/Person-Shirley-Rooney-7ceeafe5-9f9a-4315-b3da-a5b4b69c013a.json")

	person := Person{
		Thing: Thing{
			ID:        "http://api.ft.com/things/7ceeafe5-9f9a-4315-b3da-a5b4b69c013a",
			APIURL:    "http://api.ft.com/people/7ceeafe5-9f9a-4315-b3da-a5b4b69c013a",
			PrefLabel: "Shirley Rooney",
		},
		Types: []string{
			"http://www.ft.com/ontology/core/Thing",
			"http://www.ft.com/ontology/concept/Concept",
			"http://www.ft.com/ontology/person/Person",
		},
		Memberships:    []Membership{},
		DirectType:     "http://www.ft.com/ontology/person/Person",
		TwitterHandle:  "@something",
		EmailAddress:   "test@example.com",
		DescriptionXML: "Some text containing <strong>markup</strong>",
		ImageURL:       "http://someimage.jpg",
	}
	readConceptAndCompare(t, person, "7ceeafe5-9f9a-4315-b3da-a5b4b69c013a")
}

func readConceptAndCompare(t *testing.T, expected Person, uuid string) {
	actual, found, err := publicPeopleDriver.Read(uuid, testTransactionID)

	assert.NoError(t, err, "Unexpected Error occurred reading UUID: %v", uuid)
	assert.True(t, found, "Person not found with UUID: %v", uuid)
	assert.NotNil(t, actual)

	sort.Slice(expected.Memberships, func(i, j int) bool {
		return expected.Memberships[i].Organisation.ID < expected.Memberships[j].Organisation.ID
	})

	sort.Slice(actual.Memberships, func(i, j int) bool {
		return actual.Memberships[i].Organisation.ID < actual.Memberships[j].Organisation.ID
	})

	assert.Equal(t, expected.Memberships, actual.Memberships, "Expected Memberships differ from actual \nExpected: %v \nActual: %v", expected.Memberships, actual.Memberships)

	sort.Slice(expected.Labels, func(i, j int) bool {
		return expected.Labels[i] < expected.Labels[j]
	})

	sort.Slice(actual.Labels, func(i, j int) bool {
		return actual.Labels[i] < actual.Labels[j]
	})

	for _, membership := range actual.Memberships {
		sort.Slice(membership.Roles, func(i, j int) bool {
			return membership.Roles[i].ID < membership.Roles[j].ID
		})

		sort.Slice(membership.Types, func(i, j int) bool {
			return membership.Types[i] < membership.Types[j]
		})

		sort.Slice(membership.ChangeEvents, func(i, j int) bool {
			return membership.ChangeEvents[i].StartedAt < membership.ChangeEvents[j].StartedAt
		})
	}

	for _, membership := range expected.Memberships {
		sort.Slice(membership.Roles, func(i, j int) bool {
			return membership.Roles[i].ID < membership.Roles[j].ID
		})

		sort.Slice(membership.Types, func(i, j int) bool {
			return membership.Types[i] < membership.Types[j]
		})

		sort.Slice(membership.ChangeEvents, func(i, j int) bool {
			return membership.ChangeEvents[i].StartedAt < membership.ChangeEvents[j].StartedAt
		})
	}

	assert.Equal(t, expected.Labels, actual.Labels, "Expected labels differ from actual \nExpected: %v \nActual: %v", expected.Labels, actual.Labels)
	assert.Equal(t, expected.ID, actual.ID, "Expected labels differ from actual \nExpected: %v \nActual: %v", expected.ID, actual.ID)
	assert.Equal(t, expected.APIURL, actual.APIURL, "Expected API URL differ from actual \nExpected: %v \nActual: %v", expected.APIURL, actual.APIURL)
	assert.Equal(t, expected.PrefLabel, actual.PrefLabel, "Expected PrefLabel differ from actual \nExpected: %v \nActual: %v", expected.PrefLabel, actual.PrefLabel)
	assert.Equal(t, expected.BirthYear, actual.BirthYear, "Expected BirthYear differ from actual \nExpected: %v \nActual: %v", expected.BirthYear, actual.BirthYear)
	assert.Equal(t, expected.Salutation, actual.Salutation, "Expected BirthYear differ from actual \nExpected: %v \nActual: %v", expected.Salutation, actual.Salutation)
	assert.Equal(t, expected.Description, actual.Description, "Expected Description differ from actual \nExpected: %v \nActual: %v", expected.Description, actual.Description)
	assert.Equal(t, expected.DescriptionXML, actual.DescriptionXML, "Expected DescriptionXML differ from actual \nExpected: %v \nActual: %v", expected.DescriptionXML, actual.DescriptionXML)
	assert.Equal(t, expected.ImageURL, actual.ImageURL, "Expected ImageURL differ from actual \nExpected: %v \nActual: %v", expected.ImageURL, actual.ImageURL)
	assert.Equal(t, expected.EmailAddress, actual.EmailAddress, "Expected EmailAddress differ from actual \nExpected: %v \nActual: %v", expected.EmailAddress, actual.EmailAddress)
	assert.Equal(t, expected.TwitterHandle, actual.TwitterHandle, "Expected TwitterHandle differ from actual \nExpected: %v \nActual: %v", expected.TwitterHandle, actual.TwitterHandle)
	assert.Equal(t, expected.FacebookProfile, actual.FacebookProfile, "Expected FacebookProfile differ from actual \nExpected: %v \nActual: %v", expected.FacebookProfile, actual.FacebookProfile)

	sort.Slice(expected.Types, func(i, j int) bool {
		return expected.Types[i] < expected.Types[j]
	})

	sort.Slice(actual.Types, func(i, j int) bool {
		return actual.Types[i] < actual.Types[j]
	})

	assert.Equal(t, expected.Types, actual.Types, "Expected Types differ from actual \nExpected: %v \nActual: %v", expected.Types, actual.Types)
	assert.Equal(t, expected.DirectType, actual.DirectType, "Expected DirectType differ from actual \nExpected: %v \nActual: %v", expected.DirectType, actual.DirectType)

	assert.Equal(t, expected.Labels, actual.Labels, "Expected Labels differ from actual \nExpected: %v \nActual: %v", expected.Labels, actual.Labels)
	assert.True(t, reflect.DeepEqual(expected, actual), "Actual person differs from expected: Expected: %v, Actual: %v", expected, actual)
}

func readJSONtoPerson(t *testing.T, pathToJSONFile string) Person {
	f, err := os.Open(pathToJSONFile)
	assert.NoError(t, err)
	dec := json.NewDecoder(f)
	p := Person{}
	err = dec.Decode(&p)
	assert.NoError(t, err)
	return p
}

func writeJSONToConceptsService(t *testing.T, pathToJSONFile string) {
	f, err := os.Open(pathToJSONFile)
	assert.NoError(t, err)
	dec := json.NewDecoder(f)
	inst, _, errr := conceptsDriver.DecodeJSON(dec)
	assert.NoError(t, errr)

	_, errs := conceptsDriver.Write(inst, "TRANS_ID")
	assert.NoError(t, errs)
}

func writeJSONToService(t *testing.T, service baseftrwapp.Service, pathToJSONFile string) {
	f, err := os.Open(pathToJSONFile)
	assert.NoError(t, err)
	dec := json.NewDecoder(f)
	inst, _, errr := service.DecodeJSON(dec)
	assert.NoError(t, errr)

	errs := service.Write(inst, "TRANS_ID")
	assert.NoError(t, errs)
}

func cleanDB(db neoutils.NeoConnection, t *testing.T, uuids ...string) {
	qs := make([]*neoism.CypherQuery, len(uuids))
	for i, uuid := range uuids {
		qs[i] = &neoism.CypherQuery{
			Statement: fmt.Sprintf(`
			MATCH (a:Thing {uuid: "%s"})
			OPTIONAL MATCH (a)<-[ii:IDENTIFIES]-(i)
			OPTIONAL MATCH (a)-[annotation]-(c:Content)
			OPTIONAL MATCH (a)-[eq:EQUIVALENT_TO]-(canonical)
			OPTIONAL MATCH (canonical)<-[eq2:EQUIVALENT_TO]-(concepts)
			DETACH DELETE ii, i, annotation, eq, eq2, canonical, a`, uuid)}
	}
	err := db.CypherBatch(qs)
	assert.NoError(t, err, fmt.Sprintf("Error executing clean up cypher. Error: %v", err))
}
