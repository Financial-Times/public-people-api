package people

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNeoReadStructToPersonMandatoryFields checks that madatory fields are set even if they are empty or nil / null
func TestNeoReadStructToPersonMandatoryFields(t *testing.T) {
	assert := assert.New(t)
	expected := `{"id":"http://api.ft.com/things/","apiUrl":"","types":null,"memberships":[]}`
	person := neoReadStructToPerson(neoReadStruct{})
	personJSON, err := json.Marshal(person)
	assert.NoError(err, "Unable to marshal Person to JSON")
	assert.Equal(expected, string(personJSON), "Actual: %s doesn't match Expected: %s", string(personJSON), expected)
}

func TestNeoReadStructToPersonMultipleMemberships(t *testing.T) {
	assert := assert.New(t)
	expected := `{"id":"http://api.ft.com/things/","apiUrl":"","types":null,"memberships":[]}`
	neoStruct := neoReadStruct{}
	// neoStruct.M = Membership[2]
	person := neoReadStructToPerson(neoStruct)
	personJSON, err := json.Marshal(person)
	assert.NoError(err, "Unable to marshal Person to JSON")
	assert.Equal(expected, string(personJSON), "Actual: %s doesn't match Expected: %s", string(personJSON), expected)
}
