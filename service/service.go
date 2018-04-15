package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/sbogacz/wouldyoutatter/contender"
	"github.com/sbogacz/wouldyoutatter/dynamostore"

	log "github.com/sirupsen/logrus"
)

// Service holds the necessary clients to run the wouldyoutatter
// service
type Service struct {
	config         Config
	contenderStore *contender.Store
	router         *chi.Mux
	cancel         chan struct{}
}

// New tries to cerate a new instance of Service
func New(c Config, storer dynamostore.Storer) *Service {
	return &Service{
		config:         c,
		contenderStore: contender.NewStore(storer),
		router:         chi.NewRouter(),
		cancel:         make(chan struct{}),
	}
}

// Start starts the server
func (s *Service) Start() {
	s.router.Route("/contenders", func(r chi.Router) {
		r.Post("/", s.createContender)
		r.Route("/{contenderID}", func(r chi.Router) {
			r.Get("/", s.getContender)
			r.Delete("/", s.deleteContender)
		})
	})

	h := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		Handler:      s.router,
	}

	go func() {
		<-s.cancel
		_ = h.Shutdown(context.Background())
	}()

	fmt.Printf("Listening on port: %d\n", s.config.Port)
	if err := h.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// Stop stops the server gracefully
func (s *Service) Stop() {
	s.cancel <- struct{}{}
}

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

	if err := s.contenderStore.Set(context.Background(), c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to store contender"))
		log.Errorf("failed to store contender: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Service) getContender(w http.ResponseWriter, req *http.Request) {
	contenderID := chi.URLParam(req, "contenderID")

	c, err := s.contenderStore.Get(context.Background(), contenderID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to store contender"))
		log.Errorf("failed to store contender: %v", err)
		return
	}

	if c == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("no contender found with id: %s", contenderID)))
		return
	}

	b, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to encode contender"))
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
