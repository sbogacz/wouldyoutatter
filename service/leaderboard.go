package service

import (
	"context"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s *Service) getLeaderboard(w http.ResponseWriter, req *http.Request) {
	leaderboard, err := s.leaderboard.Get(context.TODO())
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
