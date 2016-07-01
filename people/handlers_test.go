package people

import (
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/satori/go.uuid"
)

// TestNeoReadStructToPersonMandatoryFields checks that mandatory fields are set even if they are empty or nil.
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

//{
//    "apiUrl": "http://api.ft.com/people/9b3fca66-a028-468e-8a77-d8259d621ff7",
//    "id": "http://api.ft.com/things/9b3fca66-a028-468e-8a77-d8259d621ff7",
//    "labels": [
//        "Max Axilion",
//        "Max Axilion Snr",
//        "Mr Axilion"
//    ],
//    "types": [
//        "http://www.ft.com/ontology/core/Thing",
//        "http://www.ft.com/ontology/concept/Concept",
//        "http://www.ft.com/ontology/person/Person"
//    ]
//}
func TestParseInvalidID(t *testing.T) {
	assert := assert.New(t)
	expected, _ := uuid.FromString("9b3fca66-a028-468e-8a77-d8259d621ff7")
	thing := Thing{ID: "http://api.ft.com/things/9b3fca66-a028-468e-8a77-d8259d621ff7"}
	//p := Person{}
	p := Person{Thing:&thing}
	//Types: [...]string{"http://www.ft.com/ontology/core/Thing","http://www.ft.com/ontology/concept/Concept","http://www.ft.com/ontology/person/Person"},
	actual, err := extractCanonicalUUID(p)
	assert.NoError(err, "Unable to extract canonical UUID")
	assert.True(uuid.Equal(expected, actual))
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 32,
		Dial: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).Dial,
	},
}
