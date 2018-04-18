package contender

import (
	"context"
	"sort"

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

// ScoreMatchup lets you declare
func (s *MatchupStore) ScoreMatchup(ctx context.Context, winner, loser string) error {
	scoredMatchup := newScoredMatchup(winner, loser)
	return errors.Wrapf(s.db.Update(ctx, scoredMatchup), "failed to score matchup between winner %s the loser %s", winner, loser)
}

// OrderMatchup is a helper function to ensure the contenders are ordered
// lexicographically
func OrderMatchup(c1, c2 string) (contender1, contender2 string) {
	contenders := []string{c1, c2}
	sort.Strings(contenders)
	return contenders[0], contenders[1]
}

func newScoredMatchup(winner, loser string) *Matchup {
	contender1, contender2 := OrderMatchup(winner, loser)
	return &Matchup{
		Contender1:    contender1,
		Contender2:    contender2,
		contender1Won: contender1 == winner,
	}
}
