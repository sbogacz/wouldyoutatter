package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/contender"
)

// Service holds the necessary clients to run the wouldyoutatter
// service
type Service struct {
	config         Config
	contenderStore contender.Store
	router         *chi.Mux
	cancel         chan struct{}
}

// New tries to cerate a new instance of Service
func New(c Config) (*Service, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load default AWS config")
	}

	return &Service{
		config:         c,
		contenderStore: contender.NewDynamoStore(dynamodb.New(cfg)),
		router:         chi.NewRouter(),
		cancel:         make(chan struct{}),
	}, nil
}

// Start starts the server
func (s *Service) Start() {
	s.router.Route("/contender", func(r chi.Router) {
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
