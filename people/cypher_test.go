package people

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"testing"
)

// TestNeoReadStructToPersonMandatoryFields checks that madatory fields are set even if they are empty or nil / null
func TestNeoReadStructToPersonMandatoryFields(t *testing.T) {
	expected := `{"id":"http://api.ft.com/things/","apiUrl":"","types":null,"memberships":[]}`
	neo := neoReadStructToPerson(neoReadStruct{})
	personJSON, err := json.Marshal(neo)
	log.Infof("Got %s", personJSON)
	if err != nil {
		t.Errorf("Error %s", err)
	}
	if string(personJSON) != expected {
		t.Errorf("Actual: %s doesn't match Expected: %s", string(personJSON), expected)
	}
}
