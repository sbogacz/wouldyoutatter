package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sbogacz/wouldyoutatter/contender"
	log "github.com/sirupsen/logrus"
)

func (s *Service) chooseMatchup(w http.ResponseWriter, req *http.Request) {
	userID := req.Header.Get("X-Tatter-UserID")
	masterSet, err := s.masterMatchupSet.Get(context.TODO())
	if err != nil {
		http.Error(w, "failed to retrieve master matchup set", http.StatusInternalServerError)
		log.WithError(err).Error("failed to retrieve master matchup set")
		return
	}

	userSet := &contender.MatchupSet{}
	if userID != "" {
		userSet, err = s.userMatchupSet.Get(context.TODO(), userID)
		if err != nil {
			http.Error(w, "failed to retrieve user matchup set", http.StatusInternalServerError)
			log.WithError(err).Error("failed to retrieve user matchup set")
			return
		}
	}

	possibleMatchups := masterSet.Set
	seenMatchups := userSet.Set

	// if the lists are the same length, then reset the user's set
	if len(seenMatchups) == len(possibleMatchups) {
		if err := s.userMatchupSet.Delete(context.TODO(), userID); err != nil {
			http.Error(w, "failed to reset user matchup set", http.StatusInternalServerError)
			log.WithError(err).Error("failed to reset user matchup set")
			return
		}
		seenMatchups = []contender.MatchupSetEntry{}
	}

	var matchup contender.MatchupSetEntry
	// do it naively for now
	for _, possibleMatchup := range possibleMatchups {
		for _, seenMatchup := range seenMatchups {
			if seenMatchup.String() != possibleMatchup.String() {
				matchup = possibleMatchup
			}
		}
	}

	// create a token for the matchup
	token, err := s.tokenStore.CreateToken(context.TODO(), matchup.Contender1, matchup.Contender2)
	if err != nil {
		http.Error(w, "failed to create token for voting", http.StatusInternalServerError)
		log.WithError(err).Error("failed to create token for voting")
		return
	}

}

func (s *Service) getMatchupStats(w http.ResponseWriter, req *http.Request) {
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
