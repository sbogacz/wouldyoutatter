package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func (s *Service) chooseMatchup(w http.ResponseWriter, req *http.Request) {

}

func (s *Service) getMatchupStatus(w http.ResponseWriter, req *http.Request) {
	contender1, contender2 := chi.URLParam(req, "contender1"), chi.URLParam(req, "contender2")
	if contender1 == "" {
		writeErrorMsg(w, http.StatusBadRequest, "contender1 cannot be empty in order to retrieve stats")
		return
	}
	if contender2 == "" {
		writeErrorMsg(w, http.StatusBadRequest, "contender2 cannot be empty in order to retrieve stats")
		return
	}

	matchup, err := s.matchupStore.Get(context.TODO(), contender1, contender2)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	e := json.NewEncoder(w)
	if err := e.Encode(matchup); err != nil {
		log.WithError(err).Error("failed to encode matchup")
		writeErrorMsg(w, http.StatusInternalServerError, "failed to encode matchup")
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Service) voteOnMatchup(w http.ResponseWriter, req *http.Request) {
	contender1, contender2 := chi.URLParam(req, "contender1"), chi.URLParam(req, "contender2")
	if contender1 == "" {
		writeErrorMsg(w, http.StatusBadRequest, "contender1 cannot be empty in order to retrieve stats")
		return
	}
	if contender2 == "" {
		writeErrorMsg(w, http.StatusBadRequest, "contender2 cannot be empty in order to retrieve stats")
		return
	}

}
