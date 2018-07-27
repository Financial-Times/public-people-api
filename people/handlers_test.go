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
	suite.handler = NewHandler(suite.mockDriver, 0, "http://localhost:8080")
	suite.handler.RegisterHandlers(suite.router)
}

func (suite *HandlerTestSuite) TestGetPeople_Success() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	uuid := "60e54253-1e94-38df-83b1-a39804d1ac18"
	url := "http://localhost:8080/concepts/" + uuid
	fakeResponse := `{
		"id": "http://www.ft.com/thing/60e54253-1e94-38df-83b1-a39804d1ac18",
		"apiUrl": "http://api.ft.com/people/60e54253-1e94-38df-83b1-a39804d1ac18",
		"prefLabel": "Neil Cole",
		"type": "http://www.ft.com/ontology/person/Person"
	}`

	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, fakeResponse))

	person := Person{
		Thing: Thing{
			ID:        "http://api.ft.com/things/60e54253-1e94-38df-83b1-a39804d1ac18",
			APIURL:    "http://api.ft.com/people/60e54253-1e94-38df-83b1-a39804d1ac18",
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

func (suite *HandlerTestSuite) TestGetPeople_Success_CompleteResponse() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	uuid := "60e54253-1e94-38df-83b1-a39804d1ac18"
	url := "http://localhost:8080/concepts/" + uuid

	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, conceptAPICompleteResponse))

	person := Person{
		Thing: Thing{
			ID:        "http://api.ft.com/things/60e54253-1e94-38df-83b1-a39804d1ac18",
			APIURL:    "http://api.ft.com/people/60e54253-1e94-38df-83b1-a39804d1ac18",
			PrefLabel: "Neil Cole",
		},
		DirectType: "http://www.ft.com/ontology/person/Person",
		Types: []string{
			"http://www.ft.com/ontology/core/Thing",
			"http://www.ft.com/ontology/concept/Concept",
			"http://www.ft.com/ontology/person/Person",
		},
		Labels: []string{"Neil Cole"},
		Memberships: []Membership{
			Membership{
				Title: "Graduate Degree",
				Types: []string{
					"http://www.ft.com/ontology/core/Thing",
					"http://www.ft.com/ontology/concept/Concept",
					"http://www.ft.com/ontology/organisation/Membership",
				},
				DirectType: "http://www.ft.com/ontology/organisation/Membership",
				Organisation: Organisation{
					Thing: Thing{
						ID:        "http://api.ft.com/things/1d448227-8b1b-3490-aeb8-18aa699d75f8",
						APIURL:    "http://api.ft.com/organisations/1d448227-8b1b-3490-aeb8-18aa699d75f8",
						PrefLabel: "Maurice A. Deane School of Law at Hofstra University",
					},
					Types: []string{
						"http://www.ft.com/ontology/core/Thing",
						"http://www.ft.com/ontology/concept/Concept",
						"http://www.ft.com/ontology/organisation/Organisation",
					},
					DirectType: "http://www.ft.com/ontology/organisation/Organisation",
					Labels:     []string(nil),
				},
				ChangeEvents: []ChangeEvent{
					ChangeEvent{
						StartedAt: "1979-01-01",
					},
					ChangeEvent{
						EndedAt: "1982-01-01",
					},
				},
				Roles: []Role{
					Role{
						Thing: Thing{
							ID:        "http://api.ft.com/things/c89c1b9e-2bc5-3dbd-bcc5-595d2dabb4bd",
							APIURL:    "http://api.ft.com/things/c89c1b9e-2bc5-3dbd-bcc5-595d2dabb4bd",
							PrefLabel: "Graduate Degree",
						},
						Types: []string{
							"http://www.ft.com/ontology/core/Thing",
							"http://www.ft.com/ontology/concept/Concept",
							"http://www.ft.com/ontology/MembershipRole",
						},
						DirectType: "http://www.ft.com/ontology/MembershipRole",
						ChangeEvents: []ChangeEvent{
							ChangeEvent{
								StartedAt: "1979-01-01",
							},
							ChangeEvent{
								EndedAt: "1982-01-01",
							},
						},
					},
				},
			},
		},
		Salutation:      "Mr.",
		BirthYear:       1957,
		EmailAddress:    "example@example.com",
		TwitterHandle:   "@ft",
		FacebookProfile: "https://www.facebook.com/financialtimes/",
		DescriptionXML:  "foobar",
		ImageURL:        "https://www.ft.com/__origami/service/image/v2/images/raw/fthead-v1:merryn-somerset-webb?source=next",
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

// Test case for the concept API response type is not people
func (suite *HandlerTestSuite) TestGetPeople_NotFound_NoPerson() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	uuid := "2d3e16e0-61cb-4322-8aff-3b01c59f4daa"
	url := "http://localhost:8080/concepts/" + uuid
	fakeResponse := `{
		"id": "http://www.ft.com/thing/60e54253-1e94-38df-83b1-a39804d1ac18",
		"apiUrl": "http://api.ft.com/people/60e54253-1e94-38df-83b1-a39804d1ac18",
		"prefLabel": "Brand",
		"type": "http://www.ft.com/ontology/product/Brand"
	}`

	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, fakeResponse))

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
		"id": "http://www.ft.com/thing/2d3e16e0-61cb-4322-8aff-3b01c59f4daa",
		"apiUrl": "http://api.ft.com/people/2d3e16e0-61cb-4322-8aff-3b01c59f4daa",
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

var conceptAPICompleteResponse = `{
  "id": "http://www.ft.com/thing/60e54253-1e94-38df-83b1-a39804d1ac18",
  "apiUrl": "http://api.ft.com/people/60e54253-1e94-38df-83b1-a39804d1ac18",
  "type": "http://www.ft.com/ontology/person/Person",
	"prefLabel": "Neil Cole",
	"descriptionXML": "foobar",
	"imageURL": "https://www.ft.com/__origami/service/image/v2/images/raw/fthead-v1:merryn-somerset-webb?source=next",
  "alternativeLabels": [
    {
      "type": "http://www.ft.com/ontology/Alias",
      "value": "Neil Cole"
    }
	],
	"account": [
		{
			"type": "http://www.ft.com/ontology/emailAddress",
			"value": "example@example.com"
		},
		{
			"type": "http://www.ft.com/ontology/twitterHandle",
			"value": "@ft"
		},
		{
			"type": "http://www.ft.com/ontology/facebookProfile",
			"value": "https://www.facebook.com/financialtimes/"
		}
	],
  "salutation": "Mr.",
  "birthYear": 1957,
  "relatedConcepts": [
    {
      "concept": {
        "id": "http://www.ft.com/thing/ea3e354e-13dc-3287-8950-230f3c6416d0",
        "apiUrl": "http://api.ft.com/concepts/ea3e354e-13dc-3287-8950-230f3c6416d0",
        "type": "http://www.ft.com/ontology/organisation/Membership",
        "prefLabel": "Graduate Degree",
        "alternativeLabels": [
          {
            "type": "http://www.ft.com/ontology/Alias",
            "value": "Graduate Degree"
          }
        ],
        "changeEvents": [
          {
            "startedAt": "1979-01-01"
          },
          {
            "endedAt": "1982-01-01"
          }
        ],
        "relatedConcepts": [
          {
            "concept": {
              "id": "http://www.ft.com/thing/1d448227-8b1b-3490-aeb8-18aa699d75f8",
              "apiUrl": "http://api.ft.com/concepts/1d448227-8b1b-3490-aeb8-18aa699d75f8",
              "type": "http://www.ft.com/ontology/organisation/Organisation",
              "prefLabel": "Maurice A. Deane School of Law at Hofstra University",
              "alternativeLabels": [],
              "countryOfIncorporation": "US"
            },
            "predicate": "http://www.ft.com/ontology/membershipOrganisation"
          },
          {
            "concept": {
              "id": "http://www.ft.com/thing/c89c1b9e-2bc5-3dbd-bcc5-595d2dabb4bd",
              "apiUrl": "http://api.ft.com/concepts/c89c1b9e-2bc5-3dbd-bcc5-595d2dabb4bd",
              "type": "http://www.ft.com/ontology/MembershipRole",
              "prefLabel": "Graduate Degree",
              "alternativeLabels": [
                {
                  "type": "http://www.ft.com/ontology/Alias",
                  "value": "Graduate Degree"
                }
              ],
              "changeEvents": [
                {
                  "startedAt": "1979-01-01"
                },
                {
                  "endedAt": "1982-01-01"
                }
              ]
            },
            "predicate": "http://www.ft.com/ontology/membershipRole"
          }
        ]
      },
      "predicate": "http://www.ft.com/ontology/membership"
    }
  ]
}`
