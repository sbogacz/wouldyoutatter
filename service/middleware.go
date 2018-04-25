package service

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
	log "github.com/sirupsen/logrus"
)

func (s *Service) checkMasterKey(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.WithField("headers", req.Header).Debug("checking")
		key := req.Header.Get("X-Tatter-Master")
		if key == "" {
			log.Debug("no key")
			http.Error(w, "missing key for desired operations", http.StatusUnauthorized)
			return
		}

		if key != s.config.MasterKey {
			log.Debug("wrong key")
			http.Error(w, "wrong key for desired operation", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, req)
	})
}

func (s *Service) validateToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("X-Tatter-Token")
		if token == "" {
			log.Debug("no token")
			http.Error(w, "missing token for voting", http.StatusUnauthorized)
			return
		}

		contender1, contender2 := chi.URLParam(req, "contender1"), chi.URLParam(req, "contender2")
		ok, err := s.tokenStore.ValidateToken(req.Context(), token, contender1, contender2)
		if err != nil {
			if dynamostore.NotFoundError(err) {
				http.Error(w, "invalid token for voting", http.StatusUnauthorized)
				return
			}
			http.Error(w, "failed to authenticate token", http.StatusInternalServerError)
			log.WithError(err).Error("failed to authenticate token against the database")
			return
		}
		if !ok {
			http.Error(w, "invalid token for voting", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, req)
	})
}
