package people

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Financial-Times/go-logger"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"gopkg.in/jarcoal/httpmock.v1"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server    *httptest.Server
	personURL string
	isFound   bool
)

type HandlerTestSuite struct {
	suite.Suite
	mockDriver *MockDriver
	router     *mux.Router
	handler    *Handler
}

func (suite *HandlerTestSuite) SetupTest() {
	logger.InitDefaultLogger("handler-test")
	suite.router = mux.NewRouter()
	suite.mockDriver = &MockDriver{}
	suite.handler = NewHandler(suite.mockDriver, 0, "http://localhost:8080/concepts")
	suite.handler.RegisterHandlers(suite.router)
}

func (suite *HandlerTestSuite) TestGetPeople_Success() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	uuid := "60e54253-1e94-38df-83b1-a39804d1ac18"
	url := "http://localhost:8080/concepts/" + uuid
	fakeResponse := `{
		"id": "http://api.ft.com/things/60e54253-1e94-38df-83b1-a39804d1ac18",
		"apiUrl": "http://api.ft.com/concepts/60e54253-1e94-38df-83b1-a39804d1ac18",
		"prefLabel": "Neil Cole",
		"type": "http://www.ft.com/ontology/person/Person"
	}`

	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, string(fakeResponse)))

	person := Person{
		Thing: Thing{
			ID:        "http://api.ft.com/things/60e54253-1e94-38df-83b1-a39804d1ac18",
			APIURL:    "http://api.ft.com/concepts/60e54253-1e94-38df-83b1-a39804d1ac18",
			PrefLabel: "Neil Cole",
		},
		DirectType: "http://www.ft.com/ontology/person/Person",
		Types: []string{
			"http://www.ft.com/ontology/core/Thing",
			"http://www.ft.com/ontology/concept/Concept",
			"http://www.ft.com/ontology/person/Person",
		},
	}

	req := newRequest("GET", "/people/"+uuid, "")
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	retPerson := Person{}
	json.NewDecoder(rec.Result().Body).Decode(&retPerson)
	suite.Equal(http.StatusOK, rec.Result().StatusCode)
	suite.Equal(person, retPerson)
}

func (suite *HandlerTestSuite) TestGetPeople_NotFound() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	uuid := "2d3e16e0-61cb-4322-8aff-3b01c59f4daa"
	url := "http://localhost:8080/concepts/" + uuid
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(404, "Not found"))

	req := newRequest("GET", "/people/"+uuid, "")
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	msg := &errMsg{
		Message: personNotFoundMsg,
	}
	returnMsg := &errMsg{}
	json.NewDecoder(rec.Result().Body).Decode(returnMsg)
	suite.Equal(msg, returnMsg)
	suite.Equal(http.StatusNotFound, rec.Result().StatusCode)
}

func (suite *HandlerTestSuite) TestGetPeople_BadRequest() {
	req := newRequest("GET", "/people/BOO", "")
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	msg := &errMsg{
		Message: badRequestMsg,
	}

	returnMsg := &errMsg{}
	json.NewDecoder(rec.Result().Body).Decode(returnMsg)
	suite.Equal(msg, returnMsg)
	suite.Equal(http.StatusBadRequest, rec.Result().StatusCode)
}

func (suite *HandlerTestSuite) TestGetPeople_Redirect() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	uuid := "70f4732b-7f7d-30a1-9c29-0cceec23760e"
	canonicalUUID := "2d3e16e0-61cb-4322-8aff-3b01c59f4daa"

	url := "http://localhost:8080/concepts/" + uuid
	fakeResponse := `{
		"id": "http://api.ft.com/things/2d3e16e0-61cb-4322-8aff-3b01c59f4daa",
		"apiUrl": "http://api.ft.com/concepts/2d3e16e0-61cb-4322-8aff-3b01c59f4daa",
		"prefLabel": "Someone",
		"type": "http://www.ft.com/ontology/person/Person"
	}`

	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, fakeResponse))

	req := newRequest("GET", "/people/"+uuid, "")
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	msg := &errMsg{
		Message: fmt.Sprintf(redirectedPerson, uuid, canonicalUUID),
	}

	returnMsg := &errMsg{}
	json.NewDecoder(rec.Result().Body).Decode(returnMsg)
	suite.Equal(msg, returnMsg)
	suite.Equal(http.StatusMovedPermanently, rec.Result().StatusCode)
}

func (suite *HandlerTestSuite) TestGetPeople_InternalError() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	uuid := "70f4732b-7f7d-30a1-9c29-0cceec23760e"
	httpmock.RegisterResponder("GET", "http://localhost:8080/people/"+uuid, httpmock.NewStringResponder(500, string("Some error")))

	req := newRequest("GET", "/people/"+uuid, "")
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	msg := &errMsg{
		Message: personUnableToBeRetrieved,
	}

	returnMsg := &errMsg{}
	json.NewDecoder(rec.Result().Body).Decode(returnMsg)
	suite.Equal(msg, returnMsg)
	suite.Equal(http.StatusInternalServerError, rec.Result().StatusCode)
}

func (suite *HandlerTestSuite) TestGetPeople_MethodNotAllowedOnPost() {
	uuid := "70f4732b-7f7d-30a1-9c29-0cceec23760e"
	req := newRequest("POST", "/people/"+uuid, "")
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)
	suite.Equal(http.StatusMethodNotAllowed, rec.Result().StatusCode)
	suite.mockDriver.AssertExpectations(suite.T())

}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func newRequest(method, url string, body string) *http.Request {
	var payload io.Reader
	if body != "" {
		payload = bytes.NewReader([]byte(body))
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		panic(err)
	}
	return req
}

type errMsg struct {
	Message string `json:"message"`
}
