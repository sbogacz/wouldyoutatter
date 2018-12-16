package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/sbogacz/wouldyoutatter/contender"
	log "github.com/sirupsen/logrus"
)

const (
	// CookieKey refers to the name for the cookie we create
	CookieKey = "wouldyoutatterID"
)

// MatchupResp is the response we'll use for our matchups/random endpoint
type MatchupResp struct {
	Contender1 contender.Contender `json:"contender_1"`
	Contender2 contender.Contender `json:"contender_2"`
	VoteURL    string              `json:"vote_url"` // we don't record this in the DB, but we use it in the API
	remove     bool
}

// VotePayload is the struct of the expected payload on vote POSTs
type VotePayload struct {
	Winner string `json:"winner"`
}

func (s *Service) chooseMatchup(w http.ResponseWriter, req *http.Request) {
	userIDCookie, err := req.Cookie(CookieKey)
	var userID string
	var newUser bool
	// if we saw an error, that's because the cookie wasn't found
	if err != nil {
		uid, err := uuid.NewV4()
		if err != nil {
			log.WithError(err).Error("couldn't generate a user ID to put in the cookie")
		} else {
			userID = uid.String()
		}
		// put the token in the response
		http.SetCookie(w, &http.Cookie{
			Name:  CookieKey,
			Value: userID,
		})

		newUser = true
	} else {
		userID = userIDCookie.Value
	}

	log.WithField("userID", userID).Debug("getting new matchup")
	masterSet, err := s.masterMatchupSet.Get(context.TODO())
	if err != nil {
		http.Error(w, "failed to retrieve master matchup set", http.StatusInternalServerError)
		log.WithError(err).Error("failed to retrieve master matchup set")
		return
	}

	userSet := &contender.MatchupSet{}
	if !newUser {
		userSet, err = s.userMatchupSet.Get(context.TODO(), userID)
		if err != nil {
			http.Error(w, "failed to retrieve user matchup set", http.StatusInternalServerError)
			log.WithError(err).Error("failed to retrieve user matchup set")
			return
		}
	}

	possibleMatchups := masterSet.Set
	seenMatchups := userSet.Set

	if len(possibleMatchups) < 1 {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("no matchups currently available"))
		return
	}

	// if the lists are the same length, then reset the user's set
	if len(seenMatchups) == len(possibleMatchups) {
		if deleteErr := s.userMatchupSet.Delete(context.TODO(), userID); deleteErr != nil {
			http.Error(w, "failed to reset user matchup set", http.StatusInternalServerError)
			log.WithError(deleteErr).Error("failed to reset user matchup set")
			return
		}
		seenMatchups = []contender.MatchupSetEntry{}
	}

	matchup := chooseNewMatchup(possibleMatchups, seenMatchups)

	// create a token for the matchup
	token, err := s.tokenStore.CreateToken(context.TODO(), matchup.Contender1, matchup.Contender2)
	if err != nil {
		http.Error(w, "failed to create token for voting", http.StatusInternalServerError)
		log.WithError(err).Error("failed to create token for voting")
		return
	}

	// get the rest of the contender's data for the client
	contender1, err := s.contenderStore.Get(context.TODO(), matchup.Contender1)
	if err != nil {
		log.WithError(err).Error("failed to retrieve contender1 data for matchup")
		http.Error(w, "failed to retrieve matchup", http.StatusInternalServerError)
		return
	}

	contender2, err := s.contenderStore.Get(context.TODO(), matchup.Contender2)
	if err != nil {
		log.WithError(err).Error("failed to retrieve contender2 data for matchup")
		http.Error(w, "failed to retrieve matchup", http.StatusInternalServerError)
		return
	}

	// mark this matchup as shown to the user
	if err := s.userMatchupSet.Add(context.TODO(), userID, matchup.Contender1, matchup.Contender2); err != nil {
		log.WithError(err).Error("failed to record seen matchup")
	}

	newURLBase := strings.Split(req.URL.String(), "/random")[0]
	newURL := fmt.Sprintf("%s/%s/%s/vote?token=%s", newURLBase, matchup.Contender1, matchup.Contender2, token.ID)

	resp := &MatchupResp{
		Contender1: *contender1,
		Contender2: *contender2,
		VoteURL:    newURL,
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		log.WithError(err).Error("failed to marshal matchup response")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (s *Service) getMatchupStats(w http.ResponseWriter, req *http.Request) {
	contender1, contender2 := chi.URLParam(req, "contender1"), chi.URLParam(req, "contender2")
	if contender1 == "" {
		http.Error(w, "contender1 cannot be empty in order to retrieve stats", http.StatusBadRequest)
		return
	}
	if contender2 == "" {
		http.Error(w, "contender2 cannot be empty in order to retrieve stats", http.StatusBadRequest)
		return
	}

	// persist the cookie since we might be getting redirected here
	userIDCookie, err := req.Cookie(CookieKey)
	// if we saw an error, that's because the cookie wasn't found
	if err == nil {
		// put the token in the request since we redirect
		http.SetCookie(w, userIDCookie)
	}

	matchup, err := s.matchupStore.Get(context.TODO(), contender1, contender2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e := json.NewEncoder(w)
	if err := e.Encode(matchup); err != nil {
		log.WithError(err).Error("failed to encode matchup")
		http.Error(w, "failed to encode matchup", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Service) voteOnMatchup(w http.ResponseWriter, req *http.Request) {
	contender1, contender2 := chi.URLParam(req, "contenderID1"), chi.URLParam(req, "contenderID2")
	if contender1 == "" {
		http.Error(w, "contender1 cannot be empty in order to vote", http.StatusBadRequest)
		return
	}
	if contender2 == "" {
		http.Error(w, "contender2 cannot be empty in order to vote", http.StatusBadRequest)
		return
	}

	// validate token
	token := req.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "must provide a valid token in order to vote", http.StatusUnauthorized)
		return
	}
	ok, err := s.tokenStore.ValidateToken(context.TODO(), token, contender1, contender2)
	if err != nil {
		http.Error(w, "failed to validate token", http.StatusInternalServerError)
		log.WithError(err).Error("failed to authenticate token against Dynamo")
		return
	}

	if !ok {
		http.Error(w, "token not valid for the matchup", http.StatusUnauthorized)
		return
	}

	v := VotePayload{}
	d := json.NewDecoder(req.Body)
	defer req.Body.Close()

	if err := d.Decode(&v); err != nil {
		http.Error(w, "couldn't decode vote payload", http.StatusBadRequest)
		return
	}
	if v.Winner != contender1 && v.Winner != contender2 {
		http.Error(w, "can only vote for a winner within the matchup", http.StatusBadRequest)
		return
	}

	loser := contender2
	if v.Winner == contender2 {
		loser = contender1
	}
	// so the token is valid, now VOTE!
	if err := s.matchupStore.ScoreMatchup(context.TODO(), v.Winner, loser); err != nil {
		http.Error(w, "failed to record vote", http.StatusInternalServerError)
		log.WithError(err).Error("failed to record vote in DB")
		return
	}

	// update the contender table
	if err := s.contenderStore.DeclareWinner(context.TODO(), v.Winner); err != nil {
		http.Error(w, "failed to record vote", http.StatusInternalServerError)
		log.WithError(err).Error("failed to update contender store with winner")
		return
	}

	if err := s.contenderStore.DeclareLoser(context.TODO(), loser); err != nil {
		http.Error(w, "failed to record vote", http.StatusInternalServerError)
		log.WithError(err).Error("failed to update contender store with loser")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func chooseNewMatchup(possibleMatchups []contender.MatchupSetEntry, seenMatchups []contender.MatchupSetEntry) contender.MatchupSetEntry {
	// if we haven't seen any, choose a random one
	if len(seenMatchups) == 0 {
		rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
		return possibleMatchups[rand.Intn(len(possibleMatchups))]
	}

	// if we've seen some, shuffle
	shuffledMatchups := make([]contender.MatchupSetEntry, len(possibleMatchups))
	perm := rand.Perm(len(possibleMatchups))
	for i, v := range perm {
		shuffledMatchups[v] = possibleMatchups[i]
	}

	for _, possibleMatchup := range shuffledMatchups {

		log.WithField("possible", possibleMatchup).Debug("possible matchup")
		var checkedMatchups int
		for _, seenMatchup := range seenMatchups {
			log.WithField("seen", seenMatchup).Debug("seen matchup")
			if seenMatchup.String() == possibleMatchup.String() {
				break
			}
			checkedMatchups++
		}
		if checkedMatchups != len(seenMatchups) {
			continue
		}
		return possibleMatchup
	}
	return contender.MatchupSetEntry{}
}
