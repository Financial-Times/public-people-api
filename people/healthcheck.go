package people

import (
	"net/http"
	"time"

	mux "github.com/gorilla/mux"
	metrics "github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	logger "github.com/Financial-Times/go-logger"
	httphandlers "github.com/Financial-Times/http-handlers-go/httphandlers"
	gtg "github.com/Financial-Times/service-status-go/gtg"
	st "github.com/Financial-Times/service-status-go/httphandlers"
)

type HealthcheckService struct {
	config HealthConfig
	Checks []fthealth.Check
}

type HealthConfig struct {
	AppSystemCode     string
	AppName           string
	Description       string
	ReqLoggingEnabled bool
}

func NewHealthCheckService(checks []fthealth.Check, config HealthConfig) *HealthcheckService {
	return &HealthcheckService{
		config: config,
		Checks: checks,
	}
}

func (s HealthcheckService) RegisterAdminHandlers(router *mux.Router) http.Handler {
	logger.Info("Registering admin handlers")

	timedHC := fthealth.TimedHealthCheck{
		HealthCheck: fthealth.HealthCheck{
			SystemCode:  s.config.AppSystemCode,
			Name:        s.config.AppName,
			Description: s.config.Description,
			Checks:      s.Checks,
		},
		Timeout: 10 * time.Second,
	}

	router.HandleFunc("/__health", fthealth.Handler(&timedHC))
	router.HandleFunc("/__gtg", st.NewGoodToGoHandler(s.gtg))
	router.HandleFunc(st.BuildInfoPath, st.BuildInfoHandler)

	var monitoringRouter http.Handler = router
	if s.config.ReqLoggingEnabled {
		monitoringRouter = httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), monitoringRouter)
		monitoringRouter = httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry, monitoringRouter)
	}

	return monitoringRouter
}

func (s HealthcheckService) gtg() gtg.Status {
	var sc []gtg.StatusChecker
	for _, check := range s.Checks {
		statusCheck := func() gtg.Status {
			return gtgCheck(check.Checker)
		}
		sc = append(sc, statusCheck)
	}
	return gtg.FailFastParallelCheck(sc)()
}

func gtgCheck(handler func() (string, error)) gtg.Status {
	if _, err := handler(); err != nil {
		return gtg.Status{
			GoodToGo: false,
			Message:  err.Error(),
		}
	}
	return gtg.Status{
		GoodToGo: true,
	}
}
