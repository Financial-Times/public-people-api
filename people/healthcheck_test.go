package people

import (
	"errors"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/stretchr/testify/suite"
)

type HealthCheckTestTestSuite struct {
	suite.Suite
	config  HealthConfig
	service *HealthcheckService
	checker func() (string, error)
	router  *mux.Router
}

func (suite *HealthCheckTestTestSuite) SetupTest() {
	suite.router = mux.NewRouter()
	suite.config = HealthConfig{
		AppSystemCode:     "appSystemCode",
		AppName:           "appName",
		Description:       "appDescription",
		ReqLoggingEnabled: false,
	}
	expStatus := "Generic error"
	expError := errors.New("generic.error")
	checked := false
	suite.checker = func() (string, error) {
		checked = true
		return expStatus, expError
	}
	checks := []fthealth.Check{
		{
			BusinessImpact:   "Tests service checker",
			Name:             "Test healthcheck",
			PanicGuide:       "https://dewey.in.ft.com/view/system/public-people-api",
			Severity:         1,
			TechnicalSummary: "Does nothing",
			Checker:          suite.checker,
		},
	}
	suite.service = NewHealthCheckService(checks, suite.config)
	suite.service.RegisterAdminHandlers(suite.router)
	status, err := suite.service.Checks[0].Checker()
	suite.Equal(expStatus, status)
	suite.Equal(expError, err)
	suite.True(checked)
}

func (suite *HealthCheckTestTestSuite) TestGtg_Success() {
	req := newRequest("GET", "/__gtg", "")
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Equal(http.StatusServiceUnavailable, rec.Result().StatusCode)
}

func TestHealthCheckTestSuite(t *testing.T) {
	suite.Run(t, new(HealthCheckTestTestSuite))
}
