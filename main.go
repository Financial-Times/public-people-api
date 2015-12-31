package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Financial-Times/go-fthealth/v1a"
	"github.com/Financial-Times/neoism"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
)

var peopleDriver PeopleDriver
var membershipDriver MembershipDriver

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
	db, err := neoism.Connect(neoURL)
	if err != nil {
		panic(err)
	}
	peopleDriver = NewPeopleCypherDriver(db)
	membershipDriver = NewMembershipCypherDriver(db)
	r := mux.NewRouter()

	// Healthchecks and standards first
	r.HandleFunc("/__health", v1a.Handler("PeopleReadWriteNeo4j Healthchecks",
		"Checks for accessing neo4j", healthCheck()))
	r.HandleFunc("/ping", ping)

	// Then API specific ones:
	// TODO wonder if we should use a regex here since this won't match /people or /people/
	r.HandleFunc("/people/{uuid}", getPerson).Methods("GET")

	if err := http.ListenAndServe(":"+port, handlers.CombinedLoggingHandler(os.Stdout, r)); err != nil {
		log.Printf("Unable to start server: %v\n", err)
		panic(err)
	}
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
	person := peopleDriver.Read(uuid)
	memberships, _, err := membershipDriver.FindMembershipsByPersonUUID(uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	person["memberships"] = memberships

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	log.Printf("\n\nPerson with memberships %+v\n", person)
	json.NewEncoder(w).Encode(person)
}
