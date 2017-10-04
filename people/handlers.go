package people

import (
	"encoding/json"
	"net/http"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"strings"
	"regexp"
	"github.com/Financial-Times/transactionid-utils-go"
	"fmt"
)

const (
	urlPrefix = "http://api.ft.com/things/"
	validUUID = "([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$"
)

// PeopleDriver for cypher queries
var PeopleDriver Driver
var CacheControlHeader string

//var maxAge = 24 * time.Hour

// HealthCheck does something
func HealthCheck() fthealth.Check {
	return fthealth.Check{
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

//GoodToGo returns a 503 if the healthcheck fails - suitable for use from varnish to check availability of a node
func GoodToGo(writer http.ResponseWriter, req *http.Request) {
	if _, err := Checker(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
}

// MethodNotAllowedHandler handles 405
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}

// GetPerson is the public API
func GetPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestId := vars["uuid"]
	transId := transactionidutils.GetTransactionIDFromRequest(r)
	w.Header().Set("X-Request-Id", transId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	validRegexp := regexp.MustCompile(validUUID)
	if requestId == "" || !validRegexp.MatchString(requestId) {
		msg := fmt.Sprintf("Invalid request id %s", requestId)
		log.WithFields(log.Fields{"UUID": requestId, "transaction_id": transId}).Error(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`"{\"message\":\"` + msg + `\"}"`))
		return
	}

	person, found, err := PeopleDriver.Read(requestId, transId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`"{\"message\": \"Person could not be retrieved\"}"`))
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`"{\"message\":\"Person ` + requestId + ` not found in DB\"}"`))
		return
	}

	canonicalId := strings.TrimPrefix(person.ID, urlPrefix)
	if strings.Compare(canonicalId, requestId) != 0 {
		log.WithFields(log.Fields{"UUID": requestId}).Info("Person " + requestId + " is concorded to " + canonicalId + "; serving redirect")
		w.Header().Set("Location", person.APIURL)
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}

	w.Header().Set("Cache-Control", CacheControlHeader)
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(person); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`"{\"message\":\"Person could not be retrieved\"}"`))
	}
}
