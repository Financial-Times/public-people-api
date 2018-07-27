package people

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"fmt"
	"github.com/Financial-Times/go-logger"
	"github.com/Financial-Times/transactionid-utils-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	urlPrefix       = "http://api.ft.com/things/"
	validUUID       = "([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$"
	contentTypeJson = "application/json; charset=UTF-8"

	personNotFoundMsg         = "Person could not be retrieved"
	personUnableToBeRetrieved = "Person could not be retrieved"
	badRequestMsg             = "Invalid UUID"
	redirectedPerson          = "Person %s is concorded to %s; serving redirect"
)

type Handler struct {
	driver               Driver
	cacheDuration        time.Duration
	publicConceptsApiURL string
}

func NewHandler(driver Driver, cacheDuration time.Duration, publicConceptsApiURL string) *Handler {
	h := &Handler{
		driver:               driver,
		cacheDuration:        cacheDuration,
		publicConceptsApiURL: publicConceptsApiURL,
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
		writeJSONStaus(w, badRequestMsg, http.StatusBadRequest)
		return
	}

	person, found, err := h.getPersonViaConceptsAPI(uuid)
	if err != nil {
		writeJSONStaus(w, personUnableToBeRetrieved, http.StatusInternalServerError)
		return
	}
	if !found {
		writeJSONStaus(w, personNotFoundMsg, http.StatusNotFound)
		return
	}

	canonicalId := strings.TrimPrefix(person.ID, urlPrefix)
	if canonicalId != uuid {
		logger.WithTransactionID(transId).WithField("UUID", uuid).Infof(redirectedPerson, uuid, canonicalId)
		redirectURL := strings.Replace(r.URL.String(), uuid, canonicalId, 1)
		w.Header().Set("Location", redirectURL)
		writeJSONStaus(w, fmt.Sprintf(redirectedPerson, uuid, canonicalId), http.StatusMovedPermanently)
		return
	}

	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s, public", strconv.FormatFloat(h.cacheDuration.Seconds(), 'f', 0, 64)))
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(person); err != nil {
		writeJSONStaus(w, "Person could not be retrieved", http.StatusInternalServerError)
	}
}

func (h *Handler) getPersonViaConceptsAPI(uuid string) (person Person, found bool, err error) {
	var p Person

	concept, err := getConcept(uuid, h.publicConceptsApiURL)
	if err != nil {
		if err.Error() == "Not found" {
			return p, false, nil
		}
		return p, false, err
	}

	if strings.Contains(concept.Type, "Person") == false {
		logger.Infof("Concept Type is not person. type %s, uuid: %s", concept.Type, uuid)
		return p, false, nil
	}

	convertToPerson(concept, &p)

	return p, true, nil
}

func getConcept(uuid string, apiURL string) (concept Concept, err error) {
	var c Concept

	u, err := url.Parse(apiURL)
	if err != nil {
		msg := fmt.Sprintf("URL of Concepts API is invalid of %s", uuid)
		logger.WithError(err).WithUUID(uuid).Error(msg)
		return c, err
	}

	u.Path = "/concepts/" + uuid
	q := u.Query()
	for _, query := range []string{"broader", "narrower", "related"} {
		q.Add("showRelationship", query)
	}
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		logger.WithError(err).Warnf("API request failed")
		return c, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return c, fmt.Errorf("Not found")
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.WithError(err).Warnf("Error reading response body")
		return c, err
	}

	if err := json.Unmarshal(bytes, &c); err != nil {
		logger.WithError(err).Warnf("Error parsing json")
		return c, err
	}
	return c, nil
}

func writeJSONStaus(rw http.ResponseWriter, message string, statusCode int) {
	rw.Header().Set("Content-Type", contentTypeJson)
	rw.WriteHeader(statusCode)
	logMsg := fmt.Sprintf(`{"message":"%s"}`, html.EscapeString(message))
	if _, err := rw.Write([]byte(logMsg)); err != nil {
		logger.WithError(err).Warnf("could not read json error")
	}
}
