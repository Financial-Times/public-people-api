package people

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO Add Test cases for more of the mapping functions and perhaps mock out back end (although ? if mocking neoism is of value)

// TestNeoReadStructToPersonMandatoryFields checks that madatory fields are set even if they are empty or nil / null
func TestNeoReadStructToPersonMandatoryFields(t *testing.T) {
	expected := `{"id":"http://api.ft.com/things/","apiUrl":"http://api.ft.com/things/","types":null,"memberships":[]}`
	person := neoReadStructToPerson(neoReadStruct{})
	personJSON, err := json.Marshal(person)
	assert := assert.New(t)
	assert.NoError(err, "Unable to marshal Person to JSON")
	assert.Equal(expected, string(personJSON))
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
