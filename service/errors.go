package service

import (
	"encoding/json"
	"net/http"
)

type errorJSON struct {
	Err string `json: error`
}

func writeError(w http.ResponseWriter, statusCode int, err error) {
	e := json.NewEncoder(w)
	if err := e.Encode(&errorJSON{Err: err.Error()}); err != nil {
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
