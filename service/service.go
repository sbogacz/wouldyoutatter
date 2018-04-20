package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
	matchupStore   *contender.MatchupStore
	router         *chi.Mux
	cancel         chan struct{}
}

// New tries to cerate a new instance of Service
func New(c Config) (*Service, error) {
	var storer dynamostore.Storer
	if c.AWSRegion == "" {
		storer = dynamostore.NewInMemoryStore()
	} else {
		cfg, err := c.AWSConfig()
		if err != nil {
			return nil, err
		}
		storer = dynamostore.New(dynamodb.New(cfg))
	}

	// set log level
	log.SetLevel(c.logLevelToLogrus())

	fmt.Printf("log level at: %s", log.GetLevel().String())
	return &Service{
		config:         c,
		contenderStore: contender.NewStore(storer),
		matchupStore:   contender.NewMatchupStore(storer),
		router:         chi.NewRouter(),
		cancel:         make(chan struct{}),
	}, nil
}

// Start starts the server
func (s *Service) Start() {
	// route the contenders endpoints
	s.router.Route("/contenders", func(r chi.Router) {
		r.Post("/", s.createContender)
		r.Route("/{contenderID}", func(r chi.Router) {
			r.Get("/", s.getContender)
			r.Delete("/", s.deleteContender)
		})
	})
	// route the matchups endpoints
	/*s.router.Route("/matchups", func(r chi.Router) {
		r.Get("/", s.chooseMatchup)
		r.Route("/{contenderID1}/{contenderID2}", func(r chi.Router) {
			r.Get("/", s.getMatchupStats)
			r.Post("/", s.voteOnMatchup)
		})
	})*/

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
	if err := h.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}

// Stop stops the server gracefully
func (s *Service) Stop() {
	s.cancel <- struct{}{}
}
