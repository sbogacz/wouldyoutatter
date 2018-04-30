package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/contender"
	"github.com/sbogacz/wouldyoutatter/dynamostore"

	log "github.com/sirupsen/logrus"
)

// Service holds the necessary clients to run the wouldyoutatter
// service
type Service struct {
	config           Config
	contenderStore   *contender.Store
	matchupStore     *contender.MatchupStore
	userMatchupSet   *contender.MatchupSetStore
	masterMatchupSet *contender.MasterMatchupSetStore
	tokenStore       *contender.TokenStore

	router *chi.Mux
	cancel chan struct{}
}

// New tries to cerate a new instance of Service
func New(c Config) (*Service, error) {
	// set log level
	log.SetLevel(c.logLevelToLogrus())
	log.SetOutput(os.Stdout)

	ret := &Service{
		config: c,
		router: chi.NewRouter(),
		cancel: make(chan struct{}),
	}
	if err := ret.configureStores(); err != nil {
		return nil, errors.Wrap(err, "failed to configure necessary stores")
	}
	// Set up very permissive CORS headers. Real use would want to
	// restrict AllowedOrigins for security.
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by browsers
	})
	ret.router.Use(corsMiddleware.Handler)

	return ret, nil
}

// Start starts the server
func (s *Service) Start() {
	// route the contenders endpoints
	s.router.Route("/contenders", func(r chi.Router) {
		r.With(s.checkMasterKey).Post("/", s.createContender)
		r.Route("/{contenderID}", func(r chi.Router) {
			r.Get("/", s.getContender)
			r.With(s.checkMasterKey).Delete("/", s.deleteContender)
		})
	})
	// route the matchups endpoints
	s.router.Route("/matchups", func(r chi.Router) {
		r.Get("/random", s.chooseMatchup)
		r.Route("/{contenderID1}/{contenderID2}", func(r chi.Router) {
			r.Get("/", s.getMatchupStats)
			r.Post("/vote", s.voteOnMatchup)
		})
	})

	// route the leaderboard
	s.router.Route("/leaderboard", func(r chi.Router) {
		r.Get("/", s.getLeaderboard)
	})
	h := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Handler:      s.router,
	}

	go func() {
		<-s.cancel
		_ = h.Shutdown(context.Background())
	}()

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

func (s *Service) configureStores() error {
	if s.config.AWSRegion == "" {
		storer := dynamostore.NewInMemoryStore()
		s.contenderStore = contender.NewStore(storer)
		s.matchupStore = contender.NewMatchupStore(storer)
		s.userMatchupSet = contender.NewMatchupSetStore(storer)
		s.masterMatchupSet = contender.NewMasterMatchupSetStore(storer)
		s.tokenStore = contender.NewTokenStore(storer)
	}
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return err
	}

	// instantiate Storers with their respective table configs
	contenderStorer := dynamostore.New(dynamodb.New(cfg), s.config.ContenderTableConfig)
	matchupStorer := dynamostore.New(dynamodb.New(cfg), s.config.MatchupTableConfig)
	userMatchupSetStorer := dynamostore.New(dynamodb.New(cfg), s.config.UserMatchupsTableConfig)
	masterMatchupSetStorer := dynamostore.New(dynamodb.New(cfg), s.config.MasterMatchupsTableConfig)
	tokenStorer := dynamostore.New(dynamodb.New(cfg), s.config.TokenTableConfig)

	// instantiate the respective stoers we need
	s.contenderStore = contender.NewStore(contenderStorer)
	s.matchupStore = contender.NewMatchupStore(matchupStorer)
	s.userMatchupSet = contender.NewMatchupSetStore(userMatchupSetStorer)
	s.masterMatchupSet = contender.NewMasterMatchupSetStore(masterMatchupSetStorer)
	s.tokenStore = contender.NewTokenStore(tokenStorer)
	return nil
}
