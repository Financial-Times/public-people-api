package people

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/go-logger"
	transactionidutils "github.com/Financial-Times/transactionid-utils-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	urlPrefix       = "http://api.ft.com/things/"
	validUUID       = "([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$"
	contentTypeJson = "application/json; charset=UTF-8"

	personNotFoundMsg         = "Person could not be retrieved"
	personUnableToBeRetrieved = "Person could not be retrieved"
	badRequestMsg             = "Invalid UUID"
	redirectedPerson          = "Person %s is concorded to %s; serving redirect"
	xPolicyHeader             = "X-Policy"
)

type Handler struct {
	cacheDuration        time.Duration
	publicConceptsApiURL string
	client               *http.Client
}

func NewHandler(cacheDuration time.Duration, publicConceptsApiURL string, c *http.Client) *Handler {
	h := &Handler{
		cacheDuration:        cacheDuration,
		publicConceptsApiURL: publicConceptsApiURL,
		client:               c,
	}
	return h
}

func (h *Handler) RegisterHandlers(router *mux.Router) {
	logger.Info("Registering handlers")
	handler := handlers.MethodHandler{
		"GET": http.HandlerFunc(h.GetPerson),
	}
	router.Handle("/people/{uuid}", handler)
}

// GetPerson is the public API
func (h *Handler) GetPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	transId := transactionidutils.GetTransactionIDFromRequest(r)
	w.Header().Set("X-Request-Id", transId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	validRegexp := regexp.MustCompile(validUUID)

	if uuid == "" || !validRegexp.MatchString(uuid) {
		logger.WithTransactionID(transId).WithField("UUID", uuid).Error(badRequestMsg)
		writeJSONStatus(w, badRequestMsg, http.StatusBadRequest)
		return
	}

	logger.Infof("Current request X-Policy header values: %s", r.Header.Get(xPolicyHeader))

	person, found, err := h.getPersonViaConceptsAPI(uuid, transId, r.Header.Get(xPolicyHeader))
	if err != nil {
		writeJSONStatus(w, personUnableToBeRetrieved, http.StatusInternalServerError)
		return
	}
	if !found {
		writeJSONStatus(w, personNotFoundMsg, http.StatusNotFound)
		return
	}

	canonicalId := strings.TrimPrefix(person.ID, urlPrefix)
	if canonicalId != uuid {
		logger.WithTransactionID(transId).WithField("UUID", uuid).Infof(redirectedPerson, uuid, canonicalId)
		redirectURL := strings.Replace(r.URL.String(), uuid, canonicalId, 1)
		w.Header().Set("Location", redirectURL)
		writeJSONStatus(w, fmt.Sprintf(redirectedPerson, uuid, canonicalId), http.StatusMovedPermanently)
		return
	}

	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s, public", strconv.FormatFloat(h.cacheDuration.Seconds(), 'f', 0, 64)))
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(person); err != nil {
		writeJSONStatus(w, "Person could not be retrieved", http.StatusInternalServerError)
	}
}

func (h *Handler) getPersonViaConceptsAPI(uuid, tid, xPolicies string) (person Person, found bool, err error) {
	var p Person

	concept, err := h.getConcept(uuid, tid, xPolicies)
	if err != nil {
		if err.Error() == "Not found" {
			return p, false, nil
		}
		return p, false, err
	}

	if strings.Contains(concept.Type, "Person") == false {
		logger.WithTransactionID(tid).Infof("Concept Type is not person. type %s, uuid: %s", concept.Type, uuid)
		return p, false, nil
	}

	convertToPerson(concept, &p)

	return p, true, nil
}

func (h *Handler) getConcept(uuid, tid, xPolicies string) (concept Concept, err error) {
	var c Concept

	u, err := url.Parse(h.publicConceptsApiURL)
	if err != nil {
		msg := fmt.Sprintf("URL of Concepts API is invalid of %s", uuid)
		logger.WithError(err).WithUUID(uuid).WithTransactionID(tid).Error(msg)
		return c, err
	}

	u.Path = "/concepts/" + uuid
	q := u.Query()
	for _, query := range []string{"related"} {
		q.Add("showRelationship", query)
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return c, err
	}
	req.Header.Set("X-Request-Id", tid)
	req.Header.Set(xPolicyHeader, xPolicies)
	resp, err := h.client.Do(req)
	if err != nil {
		logger.WithError(err).WithTransactionID(tid).Warnf("API request failed")
		return c, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return c, fmt.Errorf("Not found")
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.WithError(err).WithTransactionID(tid).Warnf("Error reading response body")
		return c, err
	}

	if err := json.Unmarshal(bytes, &c); err != nil {
		logger.WithError(err).WithTransactionID(tid).Warnf("Error parsing json")
		return c, err
	}
	return c, nil
}

func writeJSONStatus(rw http.ResponseWriter, message string, statusCode int) {
	rw.Header().Set("Content-Type", contentTypeJson)
	rw.WriteHeader(statusCode)
	logMsg := fmt.Sprintf(`{"message":"%s"}`, html.EscapeString(message))
	if _, err := rw.Write([]byte(logMsg)); err != nil {
		logger.WithError(err).Warnf("could not read json error")
	}
}

func (h *Handler) Healthchecks() fthealth.Check {
	return fthealth.Check{
		ID:               "public-concepts-api-check",
		BusinessImpact:   "Unable to respond to Public People API requests",
		Name:             "Check connectivity to public-concepts-api",
		PanicGuide:       "https://dewey.in.ft.com/runbooks/public-people-api",
		Severity:         2,
		TechnicalSummary: "Not being able to communicate with public-concepts-api means that requests for people cannot be performed.",
		Checker:          h.Checker,
	}
}

func (h *Handler) Checker() (string, error) {
	req, err := http.NewRequest("GET", h.publicConceptsApiURL+"/__gtg", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "UPP public-people-api")
	resp, err := h.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("health check returned a non-200 HTTP status: %v", resp.StatusCode)
	}
	return "Public Concepts API is healthy", nil
}
