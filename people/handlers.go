package people

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"

	"github.com/Financial-Times/go-fthealth/v1a"
	"github.com/gorilla/mux"
)

// HealthCheck does something
func HealthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact: "Unable to respond to Public People api requests",
		Checker:        Checker,
	}
}

// Checker does more stuff
func Checker() (string, error) {
	return "some message to return", nil
}

// Ping says pong
func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

// PeopleDriver for cypher queries
var PeopleDriver Driver

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
