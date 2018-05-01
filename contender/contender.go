package contender

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
)

// Contender is the model for the tattoo options
type Contender struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SVG         []byte `json:"svg"`
	Wins        int    `json:"wins"`
	Losses      int    `json:"losses"`
	Score       int    `json:"score"`
	isLoser     bool
}

// Contenders is a collection that implements Scannable
type Contenders []Contender

// NewWinner creates a new "winning" contender
func NewWinner(name string) *Contender {
	return &Contender{
		Name:    name,
		isLoser: false,
	}
}

// NewLoser creates a new "losing" contender
func NewLoser(name string) *Contender {
	return &Contender{
		Name:    name,
		isLoser: true,
	}
}

// Store uses a storer to interact with the underlying Contender db
type Store struct {
	db dynamostore.Storer
}

// NewStore takes a dynamodb Storer and uses it for the contender store
func NewStore(db dynamostore.Storer) *Store {
	return &Store{
		db: db,
	}
}

// Set lets you save a contender
func (s *Store) Set(ctx context.Context, c *Contender) error {
	return errors.Wrap(s.db.Set(ctx, c), "failed to save contender")
}

// Get lets you retrieve a contender by name
func (s *Store) Get(ctx context.Context, name string) (*Contender, error) {
	c := &Contender{Name: name}
	item, err := s.db.Get(ctx, c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve contender")

	}
	ret := item.(*Contender)
	return ret, nil
}

// Delete lets you delete a container by name
func (s *Store) Delete(ctx context.Context, name string) error {
	c := &Contender{Name: name}

	return errors.Wrap(s.db.Delete(ctx, c), "failed to delete contender")
}

// DeclareWinner lets you declarea a container a winner by name
func (s *Store) DeclareWinner(ctx context.Context, name string) error {
	winner := NewWinner(name)

	return errors.Wrapf(s.db.Update(ctx, winner), "failed to declare contender %s the winner", name)
}

// DeclareLoser lets you declarea a container a loser by name
func (s *Store) DeclareLoser(ctx context.Context, name string) error {
	loser := NewLoser(name)

	return errors.Wrapf(s.db.Update(ctx, loser), "failed to declare contender %s the loser", name)
}

// GetAll lets you retrieve all of the current contenders
func (s *Store) GetAll(ctx context.Context) (*Contenders, error) {
	cs := []Contender{}
	otherContenders := Contenders(cs)
	if err := s.db.Scan(ctx, &otherContenders); err != nil {
		return nil, errors.Wrap(err, "failed to get all contenders")
	}
	return &otherContenders, nil
}

// GetLeaderboard lets you retrieve the top N contenders
func (s *Store) GetLeaderboard(ctx context.Context, limit int) (*Contenders, error) {
	cs := []Contender{}
	leaderboard := Contenders(cs)
	if err := s.db.Query(ctx, &leaderboard, limit); err != nil {
		return nil, errors.Wrap(err, "failed to query for leaderboard")
	}
	return &leaderboard, nil
}
