package people

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server    *httptest.Server
	personURL string
	isFound   bool
)

const (
	expectedCacheControlHeader string = "special header"
)

type mockPeopleDriver struct{}

func (driver mockPeopleDriver) Read(id uuid.UUID) (person Person, found bool, err error) {
	returnPerson := Person{}
	returnPerson.Thing = &Thing{}
	returnPerson.ID = id.String()
	return returnPerson, isFound, nil
}

func (driver mockPeopleDriver) CheckConnectivity() error {
	return nil
}

func init() {
	PeopleDriver = mockPeopleDriver{}
	CacheControlHeader = expectedCacheControlHeader
	r := mux.NewRouter()
	r.HandleFunc("/people/{uuid}", GetPerson).Methods("GET")
	server = httptest.NewServer(r)
	personURL = fmt.Sprintf("%s/people", server.URL) //Grab the address for the API endpoint
	isFound = true
}

func TestHeadersOKOnFound(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	req, _ := http.NewRequest("GET", personURL+"/00000000-0000-002a-0000-00000000002a", nil)
	req.Close = true
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	assert.NoError(err)
	assert.EqualValues(200, res.StatusCode)
	assert.Equal(expectedCacheControlHeader, res.Header.Get("Cache-Control"))
	assert.Equal("application/json; charset=UTF-8", res.Header.Get("Content-Type"))
}

func TestReturnNotFoundIfPersonNotFound(t *testing.T) {
	assert := assert.New(t)
	isFound = false
	req, _ := http.NewRequest("GET", personURL+"/00000000-0000-002a-0000-00000000002a", nil)
	req.Close = true
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	assert.NoError(err)
	assert.EqualValues(404, res.StatusCode)
	assert.Equal("application/json; charset=UTF-8", res.Header.Get("Content-Type"))
}
