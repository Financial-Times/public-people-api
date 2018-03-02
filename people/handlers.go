package people

import (
	"encoding/json"
	"net/http"

	"fmt"
	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/service-status-go/gtg"
	"github.com/Financial-Times/transactionid-utils-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
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
		Severity:       2,
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

func GTG() gtg.Status {
	statusCheck := func() gtg.Status {
		return gtgCheck(Checker)
	}

	return gtg.FailFastParallelCheck([]gtg.StatusChecker{statusCheck})()
}

func gtgCheck(handler func() (string, error)) gtg.Status {
	if _, err := handler(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}
	return gtg.Status{GoodToGo: true}
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
		redirectURL := strings.Replace(r.URL.String(), requestId, canonicalId, 1)
		w.Header().Set("Location", redirectURL)
		w.WriteHeader(http.StatusMovedPermanently)
		w.Write([]byte(`"{\"message\":\"Person ` + requestId + ` is concorded, redirecting...\"}"`))
		return
	}

	w.Header().Set("Cache-Control", CacheControlHeader)
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(person); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`"{\"message\":\"Person could not be retrieved\"}"`))
	}
}
