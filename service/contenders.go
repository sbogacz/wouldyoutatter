package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sbogacz/wouldyoutatter/contender"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
	log "github.com/sirupsen/logrus"
)

func (s *Service) createContender(w http.ResponseWriter, req *http.Request) {
	d := json.NewDecoder(req.Body)
	defer req.Body.Close()

	c := &contender.Contender{}
	if err := d.Decode(c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to decode payload"))
		log.Errorf("failed to decode payload: %v", err)
		return
	}

	// save contender
	if err := s.contenderStore.Set(context.Background(), c); err != nil {
		http.Error(w, "failed to store contender", http.StatusInternalServerError)
		log.WithError(err).Error("failed to store contender")
		return
	}

	// get all contenders
	allContenders, err := s.contenderStore.GetAll(context.TODO())
	if err != nil {
		http.Error(w, "failed to update master matchup set", http.StatusInternalServerError)
		log.WithError(err).Error("failed to retrieve all contenders to update master matchup set")
		return
	}

	// add to master matchup set
	if err := s.masterMatchupSet.Add(context.TODO(), c.Name, allContenders); err != nil {
		http.Error(w, "failed to update master matchup set", http.StatusInternalServerError)
		log.WithError(err).Error("failed to update master matchup set")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Service) getContender(w http.ResponseWriter, req *http.Request) {
	contenderID := chi.URLParam(req, "contenderID")

	c, err := s.contenderStore.Get(context.Background(), contenderID)
	if err != nil {
		if dynamostore.NotFoundError(err) {

			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("no contender found with id: %s", contenderID)))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to retrieve contender"))
		log.Errorf("failed to retrieve contender: %v", err)
		return
	}

	if c == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("no contender found with id: %s", contenderID)))
		log.Infof("no contender found with name: %s", contenderID)
		return
	}

	b, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to encode contender"))
		log.Errorf("failed to retrieve contender: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	w.Write(b)
}

func (s *Service) deleteContender(w http.ResponseWriter, req *http.Request) {
	contenderID := chi.URLParam(req, "contenderID")

	if err := s.contenderStore.Delete(context.Background(), contenderID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to delete contender"))
		log.Errorf("failed to delete contender: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
