package main

import (
	"context"
	"net/http"
	"os"

	"fmt"

	"time"

	"github.com/Financial-Times/public-people-api/people"

	standardLog "log"
	"net"
	"os/signal"
	"syscall"

	"github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/go-logger"
	"github.com/cyberdelia/go-metrics-graphite"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/rcrowley/go-metrics"
)

const appDescription = "This service reads people from Neo4j"

func main() {
	app := cli.App("public-people-api", "A public RESTful API for accessing People in neo4j")
	appSystemCode := app.String(cli.StringOpt{
		Name:   "app-system-code",
		Value:  "public-people-api",
		Desc:   "System Code of the application",
		EnvVar: "APP_SYSTEM_CODE",
	})
	appName := app.String(cli.StringOpt{
		Name:   "app-name",
		Value:  "Public People API",
		Desc:   "Application name",
		EnvVar: "APP_NAME",
	})
	logLevel := app.String(cli.StringOpt{
		Name:   "log-level",
		Value:  "INFO",
		Desc:   "Log level to use",
		EnvVar: "LOG_LEVEL",
	})
	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "APP_PORT",
	})
	graphiteTCPAddress := app.String(cli.StringOpt{
		Name:   "graphiteTCPAddress",
		Value:  "",
		Desc:   "Graphite TCP address, e.g. graphite.ft.com:2003. Leave as default if you do NOT want to output to graphite (e.g. if running locally)",
		EnvVar: "GRAPHITE_ADDRESS",
	})
	graphitePrefix := app.String(cli.StringOpt{
		Name:   "graphitePrefix",
		Value:  "",
		Desc:   "Prefix to use. Should start with content, include the environment, and the host name. e.g. content.test.public.content.by.concept.api.ftaps59382-law1a-eu-t",
		EnvVar: "GRAPHITE_PREFIX",
	})
	logMetrics := app.Bool(cli.BoolOpt{
		Name:   "logMetrics",
		Value:  false,
		Desc:   "Whether to log metrics. Set to true if running locally and you want metrics output",
		EnvVar: "LOG_METRICS",
	})
	cacheDuration := app.String(cli.StringOpt{
		Name:   "cache-duration",
		Value:  "30s",
		Desc:   "Duration Get requests should be cached for. e.g. 2h45m would set the max-age value to '7440' seconds",
		EnvVar: "CACHE_DURATION",
	})
	requestLoggingEnabled := app.Bool(cli.BoolOpt{
		Name:   "requestLoggingEnabled",
		Value:  true,
		Desc:   "Whether to log requests",
		EnvVar: "REQUEST_LOGGING_ENABLED",
	})
	publicConceptsApiURL := app.String(cli.StringOpt{
		Name:   "publicConceptsApiURL",
		Value:  "http://localhost:8080",
		Desc:   "Public concepts API endpoint URL.",
		EnvVar: "CONCEPTS_API",
	})

	logger.InitLogger(*appSystemCode, *logLevel)
	logger.Infof("[Startup] public-people-api is starting ")

	app.Action = func() {
		logger.Infof("System code: %s, App Name: %s, Port: %s", *appSystemCode, *appName, *port)

		// This will send metrics to Graphite if a non-empty graphiteTCPAddress is passed in, or to the standard log if logMetrics is true.
		// Make sure a sensible graphitePrefix that will uniquely identify your service is passed in, e.g. "content.test.people.rw.neo4j.ftaps58938-law1a-eu-t
		if *graphiteTCPAddress != "" {
			addr, _ := net.ResolveTCPAddr("tcp", *graphiteTCPAddress)
			go graphite.Graphite(metrics.DefaultRegistry, 5*time.Second, *graphitePrefix, addr)
		}
		if *logMetrics { //useful locally
			//messy use of the 'standard' log package here as this method takes the log struct, not an interface, so can't use logrus.Logger
			go metrics.Log(metrics.DefaultRegistry, 60*time.Second, standardLog.New(os.Stdout, "metrics", standardLog.Lmicroseconds))
		}

		appConfig := people.HealthConfig{
			AppName:           *appName,
			AppSystemCode:     *appSystemCode,
			Description:       appDescription,
			ReqLoggingEnabled: *requestLoggingEnabled,
		}

		cacheDuration, durationErr := time.ParseDuration(*cacheDuration)
		if durationErr != nil {
			logger.Fatalf("Failed to parse cache duration string, %v", durationErr)
		}

		c := &http.Client{
			Transport: &http.Transport{

				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 60 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          20,
				IdleConnTimeout:       60 * time.Second,
				TLSHandshakeTimeout:   3 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConnsPerHost:   20,
			},
		}
		handler := people.NewHandler(cacheDuration, *publicConceptsApiURL, c)

		router := mux.NewRouter()
		healthCheckService := people.NewHealthCheckService([]v1_1.Check{handler.Healthchecks()}, appConfig)

		handler.RegisterHandlers(router)
		r := healthCheckService.RegisterAdminHandlers(router)

		httpServer := &http.Server{
			Addr:         fmt.Sprintf("0.0.0.0:%s", *port),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		httpServer.Handler = r

		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

		go func() {
			logger.Infof("Listening on %s", httpServer.Addr)
			if err := httpServer.ListenAndServe(); err != nil {
				logger.Errorf("HTTP server got shut down, error: %v", err)
			}
			sig <- os.Interrupt
		}()

		<-sig
		logger.Infof("Caught SIG: %#v", sig)
		logger.Infof("Wait for 5 seconds to finish processing")

		logger.Info("Shutting down HTTP server...")
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Info("HTTP server could not be properly shut down")
		}
		logger.Info("HTTP server shut down")

		time.Sleep(5 * time.Second)
		os.Exit(0)
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Errorf("App could not start, error=[%s]\n", err)
		return
	}
}
