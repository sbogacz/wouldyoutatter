package service

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (s *Service) getLeaderboard(w http.ResponseWriter, req *http.Request) {
	// let's start with 25 as a default limit
	limit := 25
	if val := req.URL.Query().Get("limit"); val != "" {
		newLimit, err := strconv.Atoi(val)
		if err != nil {
			log.WithError(err).Debug("couldn't parse provided new limit, keeping default")
		} else {
			limit = newLimit
		}
	}
	leaderboard, err := s.contenderStore.GetLeaderboard(context.TODO(), limit)
	if err != nil {
		log.WithError(err).Error("failed to retrieve leaderboard")
		http.Error(w, "failed to retrieve leaderboard", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(leaderboard)
	if err != nil {
		log.WithError(err).Error("failed to marshal leaderboard")
		http.Error(w, "failed to retrieve leaderboard", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)

}
