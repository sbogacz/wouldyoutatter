package contender

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
)

// Leaderboard is the model for the head-to-head records
// between leaderboards
type Leaderboard struct {
	Contender string
	Score     int
	Wins      int
}

// LeaderboardStore uses a storer to interact with the underlying Leaderboard db
type LeaderboardStore struct {
	db dynamostore.Storer
}

// NewLeaderboardStore takes a dynamodb Storer and uses it for the leaderboard store
func NewLeaderboardStore(db dynamostore.Storer) *LeaderboardStore {
	return &LeaderboardStore{
		db: db,
	}
}

// Set lets you save a leaderboard
func (s *LeaderboardStore) Set(ctx context.Context, m *Leaderboard) error {
	return errors.Wrap(s.db.Set(ctx, m), "failed to save leaderboard")
}

// Get lets you retrieve a leaderboard by the leaderboard names
func (s *LeaderboardStore) Get(ctx context.Context, contender string) (*Leaderboard, error) {
	m := &Leaderboard{Contender: contender}
	item, err := s.db.Get(ctx, m)
	if err != nil {
		if dynamostore.NotFoundError(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to retrieve leaderboard")

	}
	ret := item.(*Leaderboard)
	return ret, nil
}

// Delete lets you delete a container by name
func (s *LeaderboardStore) Delete(ctx context.Context, contender string) error {
	m := &Leaderboard{Contender: contender}

	return errors.Wrap(s.db.Delete(ctx, m), "failed to delete leaderboard")
}

// DeclareWinner lets you declarea a container a winner by name
func (s *LeaderboardStore) DeclareWinner(ctx context.Context, name string) error {
	winner := NewWinner(name)

	return errors.Wrapf(s.db.Update(ctx, winner), "failed to declare leaderboard %s the winner", name)
}

// DeclareLoser lets you declarea a container a loser by name
func (s *LeaderboardStore) DeclareLoser(ctx context.Context, name string) error {
	loser := NewLoser(name)

	return errors.Wrapf(s.db.Update(ctx, loser), "failed to declare leaderboard %s the loser", name)
}
