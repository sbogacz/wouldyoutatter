package contender

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
)

// Matchup is the model for the head-to-head records
// between matchups
type Matchup struct {
	Contender1     string
	Contender2     string
	Contender1Wins int
	Contender2Wins int
	contender1Won  bool
}

// MatchupStore uses a storer to interact with the underlying Matchup db
type MatchupStore struct {
	db dynamostore.Storer
}

// NewMatchupStore takes a dynamodb Storer and uses it for the matchup store
func NewMatchupStore(db dynamostore.Storer) *MatchupStore {
	return &MatchupStore{
		db: db,
	}
}

// Set lets you save a matchup
func (s *MatchupStore) Set(ctx context.Context, m *Matchup) error {
	return errors.Wrap(s.db.Set(ctx, m), "failed to save matchup")
}

// Get lets you retrieve a matchup by the matchup names
func (s *MatchupStore) Get(ctx context.Context, contender1, contender2 string) (*Matchup, error) {
	m := &Matchup{Contender1: contender1, Contender2: contender2}
	item, err := s.db.Get(ctx, m)
	if err != nil {
		if dynamostore.NotFoundError(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to retrieve matchup")

	}
	ret := item.(*Matchup)
	return ret, nil
}

// Delete lets you delete a container by name
func (s *MatchupStore) Delete(ctx context.Context, contender1, contender2 string) error {
	m := &Matchup{Contender1: contender1, Contender2: contender2}

	return errors.Wrap(s.db.Delete(ctx, m), "failed to delete matchup")
}

// DeclareWinner lets you declarea a container a winner by name
func (s *MatchupStore) DeclareWinner(ctx context.Context, name string) error {
	winner := NewWinner(name)

	return errors.Wrapf(s.db.Update(ctx, winner), "failed to declare matchup %s the winner", name)
}

// DeclareLoser lets you declarea a container a loser by name
func (s *MatchupStore) DeclareLoser(ctx context.Context, name string) error {
	loser := NewLoser(name)

	return errors.Wrapf(s.db.Update(ctx, loser), "failed to declare matchup %s the loser", name)
}
