package people

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNeoReadStructToPersonMandatoryFields checks that madatory fields are set even if they are empty or nil / null
func TestNeoReadStructToPersonMandatoryFields(t *testing.T) {
	expected := `{"id":"http://api.ft.com/things/","apiUrl":"http://api.ft.com/things/","types":null,"memberships":[]}`
	person := neoReadStructToPerson(neoReadStruct{})
	personJSON, err := json.Marshal(person)
	assert := assert.New(t)
	assert.NoError(err, "Unable to marshal Person to JSON")
	assert.Equal(expected, string(personJSON))
}

func TestTypeURIsForPeople(t *testing.T) {
	typesFromNeo := []string{"Person", "Concept", "Thing"}
	expectedURIs := []string{"http://www.ft.com/ontology/person/Person"}
	actualURIs := typeURIs(typesFromNeo)
	assert.New(t).EqualValues(expectedURIs, actualURIs)
}

func TestTypeURIsForOrganisations(t *testing.T) {
	typesFromNeo := []string{"Organisation", "Concept", "Thing"}
	expectedURIs := []string{"http://www.ft.com/ontology/organisation/Organisation"}
	actualURIs := typeURIs(typesFromNeo)
	assert.New(t).EqualValues(expectedURIs, actualURIs)
}

func TestTypeURIsForCompany(t *testing.T) {
	typesFromNeo := []string{"Organisation", "Company", "Concept", "Thing"}
	expectedURIs := []string{"http://www.ft.com/ontology/organisation/Organisation",
		"http://www.ft.com/ontology/company/Company"}
	actualURIs := typeURIs(typesFromNeo)
	assert.New(t).EqualValues(expectedURIs, actualURIs)
}

func TestTypeURIsForPublicCompany(t *testing.T) {
	typesFromNeo := []string{"PublicCompany", "Organisation", "Company", "Concept", "Thing"}
	expectedURIs := []string{"http://www.ft.com/ontology/company/PublicCompany",
		"http://www.ft.com/ontology/organisation/Organisation",
		"http://www.ft.com/ontology/company/Company"}
	actualURIs := typeURIs(typesFromNeo)
	assert.New(t).EqualValues(expectedURIs, actualURIs)
}

func TestTypeURIsForPrivateCompany(t *testing.T) {
	typesFromNeo := []string{"PrivateCompany", "Organisation", "Company", "Concept", "Thing"}
	expectedURIs := []string{"http://www.ft.com/ontology/company/PrivateCompany",
		"http://www.ft.com/ontology/organisation/Organisation",
		"http://www.ft.com/ontology/company/Company"}
	actualURIs := typeURIs(typesFromNeo)
	assert.New(t).EqualValues(expectedURIs, actualURIs)
}

func TestNeoReadStructToPersonMultipleMemberships(t *testing.T) {
	t.SkipNow()
	// Todo implement
	assert := assert.New(t)
	expected := `{"id":"http://api.ft.com/things/","apiUrl":"","types":null,"memberships":[]}`
	neoStruct := new(neoReadStruct)
	//	neoStruct.P = {ID:"111-111", Types:{"Person", "Concept", "Thing"}, PrefLabel:"Dan Murphy"}
	person := neoReadStructToPerson(*neoStruct)
	personJSON, err := json.Marshal(person)
	assert.NoError(err, "Unable to marshal Person to JSON")
	assert.Equal(expected, string(personJSON))
}
