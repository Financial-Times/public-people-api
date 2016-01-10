package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os"

	"github.com/Financial-Times/go-fthealth/v1a"
	"github.com/Financial-Times/public-people-api/people"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/jmcvetta/neoism"
)

func main() {
	log.Debugf("public-people-api %+v", os.Args)
	app := cli.App("public-people-api-neo4j", "A public RESTful API for accessing People in neo4j")
	neoURL := app.StringOpt("neo-url", "http://localhost:7474/db/data", "neo4j endpoint URL")
	port := app.StringOpt("port", "8080", "Port to listen on")

	app.Action = func() {
		runServer(*neoURL, *port)
	}

	app.Run(os.Args)
}

func runServer(neoURL string, port string) {
	db, err := neoism.Connect(neoURL)
	if err != nil {
		panic(err)
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
		log.Printf("Unable to start server: %v\n", err)
		panic(err)
	}
}
