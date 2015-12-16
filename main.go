package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

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
	// TODO wonder if we should use a regex here since this won't match /people or /people/
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
		PrefLabel: "someName",
		ID:        uuid,
		Memberships: []Membership{
			{Title: "213",
				Organisation: Organisation{
					ID: "org-123",
				},
				Roles: []Role{
					{ID: "role-123",
						ChangeEvents: []ChangeEvent{
							{StartedAt: time.Now()},
						},
					},
				},
			},
		},
	}

	// Organisation: Organisation{
	// {PrefLabel: "Marks and Sparks"},
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(person)
	w.WriteHeader(http.StatusOK)
}
