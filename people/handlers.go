package people

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"

	"github.com/Financial-Times/go-fthealth/v1a"
	"github.com/gorilla/mux"
)

// PeopleDriver for cypher queries
var PeopleDriver Driver

// HealthCheck does something
func HealthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Unable to respond to Public People api requests",
		Name:             "Check connectivity to Neo4j - neoUrl is a parameter in hieradata for this service",
		PanicGuide:       "TODO - write panic guide",
		Severity:         1,
		TechnicalSummary: "Cannot connect to Neo4j a instance with at least one person loaded in it",
		Checker:          Checker,
	}
}

// Checker does more stuff
func Checker() (string, error) {
	ok, err := PeopleDriver.CheckConnectivity()
	if ok && err == nil {
		return "Connectivity to neo4j is ok", err
	}
	return "Error connecting to neo4j", err
}

// Ping says pong
func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

// GetPerson is the public API
func GetPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if uuid == "" {
		http.Error(w, "uuid required", http.StatusBadRequest)
		return
	}
	person, found, err := PeopleDriver.Read(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	Jason, _ := json.Marshal(person)
	log.Debugf("Person(uuid:%s): %s\n", Jason)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
}
