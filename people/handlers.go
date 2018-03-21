package people

import (
	"encoding/json"
	"net/http"

	"fmt"
	"github.com/Financial-Times/transactionid-utils-go"
	"github.com/gorilla/mux"
	"regexp"
	"strings"
	"html"
	"github.com/Financial-Times/go-logger"
	"github.com/gorilla/handlers"
	"time"
	"strconv"
)

const (
	urlPrefix = "http://api.ft.com/things/"
	validUUID = "([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$"
	contentTypeJson  = "application/json; charset=UTF-8"

	personNotFoundMsg = "Person could not be retrieved"
	personUnableToBeRetrieved = "Person could not be retrieved"
	badRequestMsg = "Invalid UUID"
	redirectedPerson = "Person %s is concorded to %s; serving redirect"
)

type Handler struct {
	driver Driver
	cacheDuration time.Duration
}

func NewHandler(driver Driver, cacheDuration time.Duration) *Handler {
	h := &Handler{
		driver: driver,
		cacheDuration: cacheDuration,
	}
	return h
}

func (h *Handler) RegisterHandlers(router *mux.Router) {
	logger.Info("Registering handlers")
	handler := handlers.MethodHandler{
		"GET":    http.HandlerFunc(h.GetPerson),
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

	logger.Errorf("UUID: %s", uuid)
	if uuid == "" || !validRegexp.MatchString(uuid) {
		logger.WithTransactionID(transId).WithField("UUID", uuid).Error(badRequestMsg)
		writeJSONStaus(w, badRequestMsg, http.StatusBadRequest)
		return
	}

	person, found, err := h.driver.Read(uuid, transId)
	if err != nil {
		writeJSONStaus(w, personUnableToBeRetrieved, http.StatusInternalServerError)
		return
	}
	if !found {
		writeJSONStaus(w, personNotFoundMsg, http.StatusNotFound)
		return
	}

	canonicalId := strings.TrimPrefix(person.ID, urlPrefix)
	if strings.Compare(canonicalId, uuid) != 0 {
		logger.WithTransactionID(transId).WithField("UUID", uuid).Infof(redirectedPerson, uuid, canonicalId)
		redirectURL := strings.Replace(r.URL.String(), uuid, canonicalId, 1)
		w.Header().Set("Location", redirectURL)
		writeJSONStaus(w, fmt.Sprintf(redirectedPerson, uuid, canonicalId), http.StatusMovedPermanently)
		return
	}

	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%s, public", strconv.FormatFloat(h.cacheDuration.Seconds(), 'f', 0, 64)))
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(person); err != nil {
		writeJSONStaus(w,"Person could not be retrieved", http.StatusInternalServerError)
	}
}

func writeJSONStaus(rw http.ResponseWriter, message string, statusCode int) {
	rw.Header().Set("Content-Type", contentTypeJson)
	rw.WriteHeader(statusCode)
	logMsg := fmt.Sprintf(`{"message":"%s"}`, html.EscapeString(message))
	if _, err := rw.Write([]byte(logMsg)); err != nil {
		logger.WithError(err).Warnf("could not read json error")
	}
}

