package people

import (
	"errors"
	"testing"

	suite "github.com/stretchr/testify/suite"
	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
)

type HealthCheckTestTestSuite struct {
	suite.Suite
	config  HealthConfig
	service *HealthcheckService
	checker func() (string, error)
}

func (suite *HealthCheckTestTestSuite) SetupTest() {

}

func (suite *HealthCheckTestTestSuite) TestChecker() {
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
	suite.service = &HealthcheckService{
		config: suite.config,
		Checks: []fthealth.Check{
			fthealth.Check{
				BusinessImpact:   "Tests service checker",
				Name:             "Test healthcheck",
				PanicGuide:       "https://dewey.in.ft.com/view/system/public-people-api",
				Severity:         1,
				TechnicalSummary: "Does nothing",
				Checker:          suite.checker,
			},
		},
	}
	status, err := suite.service.Checks[0].Checker()
	suite.Equal(expStatus, status)
	suite.Equal(expError, err)
	suite.True(checked)
}

func TestHealthCheckTestSuite(t *testing.T) {
	suite.Run(t, new(HealthCheckTestTestSuite))
}
