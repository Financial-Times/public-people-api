package people

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	"github.com/Financial-Times/memberships-rw-neo4j/memberships"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/Financial-Times/organisations-rw-neo4j/organisations"
	person "github.com/Financial-Times/people-rw-neo4j/people"
	"github.com/Financial-Times/roles-rw-neo4j/roles"
	"github.com/jmcvetta/neoism"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

// TestNeoReadStructToPersonMandatoryFields checks that madatory fields are set even if they are empty or nil / null
func TestNeoReadStructToPersonMandatoryFields(t *testing.T) {
	expected := `{"id":"http://api.ft.com/things/","apiUrl":"http://api.ft.com/things/","types":null}`
	person := neoReadStructToPerson(neoReadStruct{}, "prod")
	personJSON, err := json.Marshal(person)
	assert := assert.New(t)
	assert.NoError(err, "Unable to marshal Person to JSON")
	assert.Equal(expected, string(personJSON))
}

func TestNeoReadStructToPersonEnvIsTest(t *testing.T) {
	expected := `{"id":"http://api.ft.com/things/","apiUrl":"http://test.api.ft.com/things/","types":null}`
	person := neoReadStructToPerson(neoReadStruct{}, "test")
	personJSON, err := json.Marshal(person)
	assert := assert.New(t)
	assert.NoError(err, "Unable to marshal Person to JSON")
	assert.Equal(expected, string(personJSON))
}

// uses library functions from other services to write the following objects to a local Neo instance:
// * a person called Siobhan Morden
// * 3 orgs
// * 3 memberships
// * one annotation
// * a partridge in a pear tree (maybe not)
func TestNeoReadStructToPersonIncludingMultipleMemberships(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnectionAndCheckClean(t, assert)

	peopleRW := person.NewCypherPeopleService(db)
	assert.NoError(peopleRW.Initialise())
	personId, _ := uuid.FromString("13a9d251-71db-467a-af2f-7e56a61c910a")
	writeJsonToService(peopleRW, fmt.Sprintf("./fixtures/Person-Siobhan_Morden-%s.json", personId.String()), assert)

	organisationRW := organisations.NewCypherOrganisationService(db)
	assert.NoError(organisationRW.Initialise())
	writeJsonToService(organisationRW, "./fixtures/Organisation-Parent_A-638fc0c1-c4d9-4be4-b6d9-c97a057e7d1b.json", assert)
	writeJsonToService(organisationRW, "./fixtures/Organisation-Child_Of_A-ac4be3c3-6dc1-4966-9cc5-ac824780f631.json", assert)
	writeJsonToService(organisationRW, "./fixtures/Organisation-Other-2802a267-aa96-4f68-897c-66e90d7d57e8.json", assert)

	membershipsRW := memberships.NewCypherMembershipService(db)
	assert.NoError(membershipsRW.Initialise())
	writeJsonToService(membershipsRW, "./fixtures/Membership-Siobhan_Morden-8865b295-c1f1-442e-8972-eb100dc50292.json", assert)
	writeJsonToService(membershipsRW, "./fixtures/Membership-Siobhan_Morden-d137a439-3efd-4820-9cab-c200031e3dd9.json", assert)
	writeJsonToService(membershipsRW, "./fixtures/Membership-Siobhan_Morden-e903861d-7709-4ab3-aeb4-4d272ac4d105.json", assert)

	rolesRW := roles.NewCypherDriver(db)
	assert.NoError(rolesRW.Initialise())
	writeJsonToService(rolesRW, "./fixtures/Role-0ee8e7b7-bac9-4db1-b94b-5605ce1d2907.json", assert)

	defer cleanDB(db, t, assert)
	defer membershipsRW.Delete("e903861d-7709-4ab3-aeb4-4d272ac4d105")
	defer membershipsRW.Delete("d137a439-3efd-4820-9cab-c200031e3dd9")
	defer membershipsRW.Delete("8865b295-c1f1-442e-8972-eb100dc50292")
	defer organisationRW.Delete("2802a267-aa96-4f68-897c-66e90d7d57e8")
	defer organisationRW.Delete("ac4be3c3-6dc1-4966-9cc5-ac824780f631")
	defer organisationRW.Delete("638fc0c1-c4d9-4be4-b6d9-c97a057e7d1b")
	defer rolesRW.Delete("0ee8e7b7-bac9-4db1-b94b-5605ce1d2907")
	defer peopleRW.Delete(personId.String())

	publicPeopleDriver := NewCypherDriver(db, "prod")
	person, found, err := publicPeopleDriver.Read(personId)
	assert.NoError(err)
	assert.True(found, "Person not found in database")
	assert.NotNil(person)
	assertMemberships(&person, assert)
	assert.Equal([]string{"Siobhan J Morden", "Siobhan Morden"}, *person.Labels)
	assert.Equal(fmt.Sprintf("http://api.ft.com/things/%s", personId.String()), person.ID)
	assert.Equal(fmt.Sprintf("http://api.ft.com/people/%s", personId.String()), person.APIURL)
	assert.Equal("Siobhan Morden", person.PrefLabel)
	assertListContainsAll(assert, person.Types, "http://www.ft.com/ontology/core/Thing", "http://www.ft.com/ontology/concept/Concept", "http://www.ft.com/ontology/person/Person")
	assert.Equal("http://www.ft.com/ontology/person/Person", person.DirectType)

	assert.Equal(*person.Labels, []string{"Siobhan J Morden", "Siobhan Morden"})
	assert.Equal(person.PrefLabel, "Siobhan Morden")
	assert.Equal(person.BirthYear, 1974)

	assert.Equal(person.Salutation, "Ms.")
	assert.Equal(person.Description, "Some text")
	assert.Equal(person.DescriptionXML, "Some text containing <strong>markup</strong>")
	assert.Equal(person.ImageURL, "http://someimage.jpg")
	assert.Equal(person.EmailAddress, "test@example.com")
	assert.Equal(person.TwitterHandle, "@something")
	assert.Equal(person.FacebookProfile, "the-facebook-profile")
	assert.Equal(person.LinkedinProfile, "the-linkedin-profile")
}

func TestNeoReadPersonWithCanonicalUPPID(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnectionAndCheckClean(t, assert)

	peopleRW := person.NewCypherPeopleService(db)
	assert.NoError(peopleRW.Initialise())

	personId, _ := uuid.FromString("13a9d251-71db-467a-af2f-7e56a61c910a")
	writeJsonToService(peopleRW, fmt.Sprintf("./fixtures/Person-Siobhan_Morden-%s.json", personId.String()), assert)
	defer peopleRW.Delete(personId.String())

	publicPeopleDriver := NewCypherDriver(db, "prod")
	person, found, err := publicPeopleDriver.Read(personId)
	assert.NoError(err)
	assert.True(found, "Person not found in database")
	assert.NotNil(person)

	assert.Equal(fmt.Sprintf("http://api.ft.com/things/%s", personId.String()), person.ID)
	assert.Equal(fmt.Sprintf("http://api.ft.com/people/%s", personId.String()), person.APIURL)
	assertListContainsAll(assert, person.Types, "http://www.ft.com/ontology/core/Thing", "http://www.ft.com/ontology/concept/Concept", "http://www.ft.com/ontology/person/Person")
	assert.Equal("http://www.ft.com/ontology/person/Person", person.DirectType)

}

func TestNeoReadPersonWithAlternateUPPID(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnectionAndCheckClean(t, assert)

	peopleRW := person.NewCypherPeopleService(db)
	assert.NoError(peopleRW.Initialise())

	personId, _ := uuid.FromString("13a9d251-71db-467a-af2f-7e56a61c910a")
	alternativePersonId, _ := uuid.FromString("d755c384-c302-485c-b12e-ea3c6751a6b6")
	writeJsonToService(peopleRW, fmt.Sprintf("./fixtures/Person-Siobhan_Morden-%s.json", personId.String()), assert)
	defer peopleRW.Delete(personId.String())

	publicPeopleDriver := NewCypherDriver(db, "prod")
	person, found, err := publicPeopleDriver.Read(alternativePersonId)
	person1, found, err := publicPeopleDriver.Read(personId)

	assert.NoError(err)
	assert.True(found, "Person not found in database")
	assert.NotNil(person)
	assert.Equal(person, person1)

	assert.Equal(fmt.Sprintf("http://api.ft.com/things/%s", personId.String()), person.ID)
	assert.Equal(fmt.Sprintf("http://api.ft.com/people/%s", personId.String()), person.APIURL)
	assertListContainsAll(assert, person.Types, "http://www.ft.com/ontology/core/Thing", "http://www.ft.com/ontology/concept/Concept", "http://www.ft.com/ontology/person/Person")
	assert.Equal("http://www.ft.com/ontology/person/Person", person.DirectType)

}

func assertMemberships(person *Person, assert *assert.Assertions) {
	assert.Len(person.Memberships, 3)
	organisations := make([]string, 3)
	roleIds := make(map[string]string, 1)
	for i, mem := range person.Memberships {
		organisations[i] = mem.Organisation.ID
		assert.Len(mem.Roles, 1)
		roleIds[mem.Roles[0].ID] = mem.Roles[0].PrefLabel
	}
	assert.Len(organisations, 3)
	assert.Contains(organisations, "http://api.ft.com/things/2802a267-aa96-4f68-897c-66e90d7d57e8")
	assert.Contains(organisations, "http://api.ft.com/things/ac4be3c3-6dc1-4966-9cc5-ac824780f631")
	assert.Contains(organisations, "http://api.ft.com/things/638fc0c1-c4d9-4be4-b6d9-c97a057e7d1b")

	assert.Len(roleIds, 1)
	assert.Equal(roleIds["http://api.ft.com/things/0ee8e7b7-bac9-4db1-b94b-5605ce1d2907"], "Market Strategist")
}

func writeJsonToService(service baseftrwapp.Service, pathToJsonFile string, assert *assert.Assertions) {
	f, err := os.Open(pathToJsonFile)
	assert.NoError(err)
	dec := json.NewDecoder(f)
	inst, _, errr := service.DecodeJSON(dec)
	assert.NoError(errr)
	errrr := service.Write(inst)
	assert.NoError(errrr)
}

func getDatabaseConnectionAndCheckClean(t *testing.T, assert *assert.Assertions) neoutils.NeoConnection {
	db := getDatabaseConnection(t, assert)
	cleanDB(db, t, assert)
	return db
}

func getDatabaseConnection(t *testing.T, assert *assert.Assertions) neoutils.NeoConnection {
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	conf := neoutils.DefaultConnectionConfig()
	conf.Transactional = false
	db, err := neoutils.Connect(url, conf)
	assert.NoError(err, "Failed to connect to Neo4j")
	return db
}

func cleanDB(db neoutils.NeoConnection, t *testing.T, assert *assert.Assertions) {
	qs := []*neoism.CypherQuery{
		{
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%v'}) OPTIONAL MATCH (a)-[]-(b:Identifier) DETACH DELETE a,b", "638fc0c1-c4d9-4be4-b6d9-c97a057e7d1b"),
		},
		{
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%v'}) OPTIONAL MATCH (a)-[]-(b:Identifier) DETACH DELETE a,b", "ac4be3c3-6dc1-4966-9cc5-ac824780f631"),
		},
		{
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%v'}) OPTIONAL MATCH (a)-[]-(b:Identifier) DETACH DELETE a,b", "13a9d251-71db-467a-af2f-7e56a61c910a"),
		},
		{
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%v'}) OPTIONAL MATCH (a)-[]-(b:Identifier) DETACH DELETE a,b", "0ee8e7b7-bac9-4db1-b94b-5605ce1d2907"),
		},
		{
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%v'}) OPTIONAL MATCH (a)-[]-(b:Identifier) DETACH DELETE a,b", "e903861d-7709-4ab3-aeb4-4d272ac4d105"),
		},
		{
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%v'}) OPTIONAL MATCH (a)-[]-(b:Identifier) DETACH DELETE a,b", "d137a439-3efd-4820-9cab-c200031e3dd9"),
		},
		{
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%v'}) OPTIONAL MATCH (a)-[]-(b:Identifier) DETACH DELETE a,b", "8865b295-c1f1-442e-8972-eb100dc50292"),
		},
		{
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%v'}) OPTIONAL MATCH (a)-[]-(b:Identifier) DETACH DELETE a,b", "2802a267-aa96-4f68-897c-66e90d7d57e8"),
		},
		{
			//deletes parent 'org' which only has type Thing
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%v'}) OPTIONAL MATCH (a)-[]-(b:Identifier) DETACH DELETE a,b", "7b00924d-6115-4126-9bb5-5e3cfdfc8114"),
		},
		{
			//deletes parent 'org' which only has type Thing
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%v'}) OPTIONAL MATCH (a)-[]-(b:Identifier) DETACH DELETE a,b", "195df1b2-7e04-4c70-a865-4361c71e9a6b"),
		},
	}

	err := db.CypherBatch(qs)
	assert.NoError(err)
}

func assertListContainsAll(assert *assert.Assertions, list interface{}, items ...interface{}) {
	assert.Len(list, len(items))
	for _, item := range items {
		assert.Contains(list, item)
	}
}
