package service

import (
	"context"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
)

// Token is a struct we'll leverage to control the voting part of the API
type Token struct {
	ID         string
	Contender1 string
	Contender2 string
}

// TokenStore gives us some nicer typed access to the DB
type TokenStore struct {
	db dynamostore.Storer
}

func (s *TokenStore) CreateToken(ctx context.Context, contender1, contender2 string) (*Token, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate UUID for token")
	}
	t := &Token{
		ID:         uid.String(),
		Contender1: contender1,
		Contender2: contender2,
	}
	return t, nil
}
