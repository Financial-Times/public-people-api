package people

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Financial-Times/go-fthealth/v1a"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
)

// PeopleDriver for cypher queries
var PeopleDriver Driver
var CacheControlHeader string

//var maxAge = 24 * time.Hour

// HealthCheck does something
func HealthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact: "Unable to respond to Public People api requests",
		Name:           "Check connectivity to Neo4j - neoUrl is a parameter in hieradata for this service",
		PanicGuide:     "https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/public-people-api",
		Severity:       1,
		TechnicalSummary: `Cannot connect to Neo4j. If this check fails, check that Neo4j instance is up and running. You can find
				the neoUrl as a parameter in hieradata for this service. `,
		Checker: Checker,
	}
}

// Checker does more stuff
func Checker() (string, error) {
	err := PeopleDriver.CheckConnectivity()
	if err == nil {
		return "Connectivity to neo4j is ok", err
	}
	return "Error connecting to neo4j", err
}

// Ping says pong
func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

// BuildInfoHandler - This is a stop gap and will be added to when we can define what we should display here
func BuildInfoHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "build-info")
}

// MethodNotAllowedHandler handles 405
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}

// GetPerson is the public API
func GetPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if uuid == "" {
		http.Error(w, "uuid required", http.StatusBadRequest)
		return
	}
	person, found, err := PeopleDriver.Read(uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		raven.CaptureError(err, nil)
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"Person not found."}`))
		return
	}

	w.Header().Set("Cache-Control", CacheControlHeader)
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(person); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"Person could not be retrieved, err=` + err.Error() + `"}`))
	}
}
