package people

import (
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//TestNeoReadStructToPersonMandatoryFields checks that madatory fields are set even if they are empty or nil / null
func TestCanGetAPerson(t *testing.T) {
	t.SkipNow()
	// TODO figure out how best to test handlers.
	assert := assert.New(t)
	expected := `{"id":"http://api.ft.com/things/","apiUrl":"http://api.ft.com/things/","types":null}`
	person := neoReadStructToPerson(neoReadStruct{}, "prod")
	personJSON, err := json.Marshal(person)
	assert.NoError(err, "Unable to marshal Person to JSON")
	assert.Equal(expected, string(personJSON), "Actual: %s doesn't match Expected: %s", string(personJSON), expected)
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 32,
		Dial: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).Dial,
	},
}
