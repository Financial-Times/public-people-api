package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Financial-Times/go-fthealth/v1a"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/jmcvetta/neoism"
)

var db *neoism.Database

func main() {
	fmt.Println(os.Args)
	app := cli.App("public-people-api-neo4j", "A public RESTful API for accessing People in neo4j")
	neoURL := app.StringOpt("neo-url", "http://localhost:7474/db/data", "neo4j endpoint URL")
	port := app.StringOpt("port", "8080", "Port to listen on")

	app.Action = func() {
		runServer(*neoURL, *port)
	}

	app.Run(os.Args)
}

func runServer(neoURL string, port string) {
	var err error
	db, err = neoism.Connect(neoURL)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	// Healthchecks and standards first
	r.HandleFunc("/__health", v1a.Handler("PeopleReadWriteNeo4j Healthchecks",
		"Checks for accessing neo4j", healthCheck()))
	r.HandleFunc("/ping", ping)

	// Then API specific ones:
	//r.HandleFunc("/people/{uuid}", peopleWrite).Methods("PUT")
	r.HandleFunc("/people/{uuid}", getPerson).Methods("GET")

	http.ListenAndServe(":"+port, handlers.CombinedLoggingHandler(os.Stdout, r))

}

func healthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact: "Unable to respond to Public People api requests",
		Checker:        checker,
	}
}

func checker() (string, error) {
	return "some message to return", nil
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	person := Person{
		Name: "someName",
		UUID: uuid,
	}
	json.NewEncoder(w).Encode(person)
}

// Person structure for writing to responses
type Person struct {
	Identifiers []struct {
		Authority       string `json:"authority"`
		IdentifierValue string `json:"identifierValue"`
	} `json:"identifiers"`
	Name string `json:"name"`
	UUID string `json:"uuid"`
}
