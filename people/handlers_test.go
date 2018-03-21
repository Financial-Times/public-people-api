package people

import (
	"net/http"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"github.com/Financial-Times/go-logger"
	"github.com/stretchr/testify/mock"
	"encoding/json"
	"io"
	"bytes"
	"testing"
)

var (
	server    *httptest.Server
	personURL string
	isFound   bool
)

type HandlerTestSuite struct {
	suite.Suite
	mockDriver  *MockDriver
	router      *mux.Router
	handler     *Handler
}

func (suite *HandlerTestSuite) SetupTest() {
	logger.InitDefaultLogger("handler-test")
	suite.router = mux.NewRouter()
	suite.mockDriver = &MockDriver{}
	suite.handler = NewHandler(suite.mockDriver, 0)
	suite.handler.RegisterHandlers(suite.router)
}

func (suite *HandlerTestSuite) TestGetPeople_Success() {
	uuid := "70f4732b-7f7d-30a1-9c29-0cceec23760e"

	person  := Person{
		Thing: Thing{
			ID:        "http://api.ft.com/things/70f4732b-7f7d-30a1-9c29-0cceec23760e",
			APIURL:    "http://api.ft.com/people/70f4732b-7f7d-30a1-9c29-0cceec23760e",
			PrefLabel: "Someone",
		},
	}

	suite.mockDriver.On("Read", uuid, mock.Anything).Return(person, nil)

	req := newRequest("GET", "/people/" + uuid, "")
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	retPerson := Person{}
	json.NewDecoder(rec.Result().Body).Decode(&retPerson)
	suite.Equal(http.StatusOK, rec.Result().StatusCode)
	suite.Equal(person, retPerson)
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
