package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"strings"

	"github.com/Financial-Times/go-fthealth/v1a"
	"github.com/Financial-Times/public-people-api/people"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/jmcvetta/neoism"
)

func main() {
	app := cli.App("public-people-api-neo4j", "A public RESTful API for accessing People in neo4j")
	neoURL := app.StringOpt("neo-url", "http://localhost:7474/db/data", "neo4j endpoint URL")
	port := app.StringOpt("port", "8080", "Port to listen on")
	logLevel := app.StringOpt("log-level", "INFO", "Logging level (DEBUG, INFO, WARN, ERROR)")

	app.Action = func() {
		setLogLevel(strings.ToUpper(*logLevel))
		log.Infof("public-people-api will listen on port: %s, connecting to: %s", *port, *neoURL)
		runServer(*neoURL, *port)
	}
	app.Run(os.Args)
}

func runServer(neoURL string, port string) {
	db, err := neoism.Connect(neoURL)
	if err != nil {
		log.Fatalf("Error connecting to neo4j %s", err)
	}
	people.PeopleDriver = people.NewCypherDriver(db)
	r := mux.NewRouter()

	// Healthchecks and standards first
	r.HandleFunc("/__health", v1a.Handler("PeopleReadWriteNeo4j Healthchecks",
		"Checks for accessing neo4j", people.HealthCheck()))
	r.HandleFunc("/ping", people.Ping)

	// Then API specific ones:
	// TODO wonder if we should use a regex here since this won't match /people or /people/
	r.HandleFunc("/people/{uuid}", people.GetPerson).Methods("GET")

	if err := http.ListenAndServe(":"+port, handlers.CombinedLoggingHandler(os.Stdout, r)); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}

func setLogLevel(level string) {
	switch level {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	default:
		log.Errorf("Requested log level %s is not supported, will default to INFO level", level)
		log.SetLevel(log.InfoLevel)
	}
	log.Debugf("Logging level set to %s", level)
}
