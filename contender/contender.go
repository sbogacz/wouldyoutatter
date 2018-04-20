package contender

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
)

// Contender is the model for the tattoo options
type Contender struct {
	Name        string
	Description string
	SVG         []byte
	Wins        int
	Losses      int
	Score       int
	isLoser     bool
}

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
