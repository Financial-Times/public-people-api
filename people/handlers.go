package people

import (
	"encoding/json"
	"net/http"

	"fmt"
	"github.com/Financial-Times/transactionid-utils-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	router.Handle("/people/{uuid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", handler)
}

// GetPerson is the public API
func (h *Handler) GetPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestId := vars["uuid"]
	transId := transactionidutils.GetTransactionIDFromRequest(r)
	w.Header().Set("X-Request-Id", transId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	validRegexp := regexp.MustCompile(validUUID)

	if requestId == "" || !validRegexp.MatchString(requestId) {
		msg := fmt.Sprintf("Invalid request id %s", requestId)
		log.WithFields(log.Fields{"UUID": requestId, "transaction_id": transId}).Error(msg)
		writeJSONStaus(w,msg, http.StatusInternalServerError)
		return
	}

	person, found, err := h.driver.Read(requestId, transId)
	if err != nil {
		writeJSONStaus(w,"Person could not be retrieved", http.StatusInternalServerError)
		return
	}
	if !found {
		writeJSONStaus(w,`Person ` + requestId + ` not found in DB`, http.StatusNotFound)
		return
	}

	canonicalId := strings.TrimPrefix(person.ID, urlPrefix)
	if strings.Compare(canonicalId, requestId) != 0 {
		log.WithFields(log.Fields{"UUID": requestId}).Info("Person " + requestId + " is concorded to " + canonicalId + "; serving redirect")
		redirectURL := strings.Replace(r.URL.String(), requestId, canonicalId, 1)
		w.Header().Set("Location", redirectURL)
		writeJSONStaus(w,`Person ` + requestId + ` is concorded, redirecting...`, http.StatusMovedPermanently)
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
		log.WithError(err).Warnf("could not read json error")
	}
}

