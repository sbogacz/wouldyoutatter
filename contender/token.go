package contender

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
	"github.com/urfave/cli"
)

// Token is a struct we'll leverage to control the voting part of the API
type Token struct {
	ID         string
	Contender1 string
	Contender2 string
	ExpireAt   int64
}

// TokenStore gives us some nicer typed access to the DB
type TokenStore struct {
	db dynamostore.Storer
}

// NewTokenStore takes a Storer and returns a reference to an instance
// of a token store
func NewTokenStore(db dynamostore.Storer) *TokenStore {
	return &TokenStore{
		db: db,
	}
}

// CreateToken creates a new token for the given contender combination vote
func (s *TokenStore) CreateToken(ctx context.Context, contender1, contender2 string) (*Token, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate UUID for token")
	}
	contender1, contender2 = OrderMatchup(contender1, contender2)
	t := &Token{
		ID:         uid.String(),
		Contender1: contender1,
		Contender2: contender2,
		ExpireAt:   time.Now().Add(time.Hour * 24).Unix(),
	}
	if err := s.db.Set(ctx, t); err != nil {
		return nil, errors.Wrap(err, "failed to create token")
	}
	return t, nil
}

// ValidateToken checks to see whether a given token is still valid for the given matchup
func (s *TokenStore) ValidateToken(ctx context.Context, uid, contender1, contender2 string) (bool, error) {

	item, err := s.db.Get(ctx, &Token{ID: uid})
	if err != nil {
		return false, errors.Wrap(err, "failed to validate token against the db")
	}

	// if we didn't find a matching token, mark invalid
	if item == nil {
		return false, nil
	}

	t := item.(*Token)
	// sort the inputs, to make sure we check against the right db fields
	contender1, contender2 = OrderMatchup(contender1, contender2)
	if t.Contender1 != contender1 || t.Contender2 != contender2 {
		return false, nil
	}
	return true, nil
}

// InvalidateToken is used for explicit token invalidation (like when the token is used)
func (s *TokenStore) InvalidateToken(ctx context.Context, uid string) error {
	if err := s.db.Delete(ctx, &Token{ID: uid}); err != nil {
		return errors.Wrap(err, "failed to invalidate token")
	}
	return nil
}

// TokenTableConfig allows us to set configuration details
// for the dynamo table from the app
type TokenTableConfig struct {
	TableName     string
	ReadCapacity  int
	WriteCapacity int
}

// Flags returns a slice of the configuration options for the token table
func (c *TokenTableConfig) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "token-table-name",
			EnvVar:      "TOKEN_TABLE_NAME",
			Value:       "Tokens",
			Destination: &c.TableName,
		},
		cli.IntFlag{
			Name:        "token-table-read-capacity",
			EnvVar:      "TOKEN_TABLE_READ_CAPACITY",
			Value:       5,
			Destination: &c.ReadCapacity,
		},
		cli.IntFlag{
			Name:        "token-table-write-capacity",
			EnvVar:      "TOKEN_TABLE_WRITE_CAPACITY",
			Value:       5,
			Destination: &c.WriteCapacity,
		},
	}
}
