package people

import (
	"encoding/json"
	"github.com/Financial-Times/public-people-api"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"os"
	"testing"
)

//TestNeoReadStructToPersonMandatoryFields checks that madatory fields are set even if they are empty or nil / null
func TestNeoReadStructToPersonMandatoryFields(t *testing.T) {
	// TODO figure out how best to test handlers.
	t.SkipNow()
	assert := assert.New(t)
	expected := `{"id":"http://api.ft.com/things/","apiUrl":"","types":null,"memberships":[]}`
	person := neoReadStructToPerson(neoReadStruct{})
	personJSON, err := json.Marshal(person)
	assert.NoError(err, "Unable to marshal Person to JSON")
	assert.Equal(expected, string(personJSON), "Actual: %s doesn't match Expected: %s", string(personJSON), expected)
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 32,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	},
}
