package service

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type errorJSON struct {
	Err string `json:"error"`
}

func writeError(w http.ResponseWriter, statusCode int, err error) {
	e := json.NewEncoder(w)
	if encodingErr := e.Encode(&errorJSON{Err: err.Error()}); encodingErr != nil {
		log.WithError(err).Error("failed to encode error to JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
}

func writeErrorMsg(w http.ResponseWriter, statusCode int, err string) {
	e := json.NewEncoder(w)
	if err := e.Encode(&errorJSON{Err: err}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
}
