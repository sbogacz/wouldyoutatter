package contender

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
)

const (
	userMatchupSetTableName = "UserMatchups"
	masterMatchupTableName  = "PossibleMatchups"
	masterKey               = "Master"
)

// MatchupSet is one of the possible matchup combinations
type MatchupSet struct {
	ID        string
	Set       []string
	tableName string
	entry     matchupSetEntry
}

type matchupSetEntry struct {
	Contender1 string
	Contender2 string
	remove     bool
}

func (m matchupSetEntry) String() string {
	return fmt.Sprintf("%s§%s", m.Contender1, m.Contender2)
}

// MatchupSetStore gives us some helpful methods for interacting
// with the unerlying Storer
type MatchupSetStore struct {
	db dynamostore.Storer
}

// NewMatchupSetStore takes a Storer and instantiates a MatchupSetStore
func NewMatchupSetStore(db dynamostore.Storer) *MatchupSetStore {
	return &MatchupSetStore{
		db: db,
	}
}

// AddToMatchupSet given a uid corresponding to the session, and two contenders, adds them to the set
// of matchups that uid has seen
func (s *MatchupSetStore) AddToMatchupSet(ctx context.Context, uid, contender1, contender2 string) error {
	contender1, contender2 = OrderMatchup(contender1, contender2)
	newMatchupEntry := matchupSetEntry{
		Contender1: contender1,
		Contender2: contender2,
	}

	matchupSet := &MatchupSet{
		ID:        uid,
		entry:     newMatchupEntry,
		tableName: userMatchupSetTableName,
	}
	if err := s.db.Update(ctx, matchupSet); err != nil {
		return errors.Wrapf(err, "failed to update the matchup set for ID: %s", uid)
	}
	return nil
}

// RemoveFromMatchupSet given a uid corresponding to the session, and two contenders, adds them to the set
// of matchups that uid has seen
func (s *MatchupSetStore) RemoveFromMatchupSet(ctx context.Context, uid, contender1, contender2 string) error {
	contender1, contender2 = OrderMatchup(contender1, contender2)
	newMatchupEntry := matchupSetEntry{
		Contender1: contender1,
		Contender2: contender2,
		remove:     true,
	}

	matchupSet := &MatchupSet{
		ID:        uid,
		entry:     newMatchupEntry,
		tableName: userMatchupSetTableName,
	}
	if err := s.db.Update(ctx, matchupSet); err != nil {
		return errors.Wrapf(err, "failed to update the matchup set for ID: %s", uid)
	}
	return nil
}

// RemoveMatchupSet is used to restart a matchup set when it is no longer relevant
func (s *MatchupSetStore) RemoveMatchupSet(ctx context.Context, uid string) error {
	if err := s.db.Delete(ctx, &MatchupSet{ID: uid, tableName: userMatchupSetTableName}); err != nil {
		return errors.Wrap(err, "failed to delete matchup set")
	}
	return nil
}

// MasterMatchupSetStore gives us some helpful methods for interacting
// with the unerlying Storer
type MasterMatchupSetStore struct {
	db dynamostore.Storer
}

// NewMasterMatchupSetStore takes a Storer and instantiates a MasterMatchupSetStore
func NewMasterMatchupSetStore(db dynamostore.Storer) *MasterMatchupSetStore {
	return &MasterMatchupSetStore{
		db: db,
	}
}

// AddToMatchupSet given a uid corresponding to the session, and two contenders, adds them to the set
// of matchups that uid has seen
func (s *MasterMatchupSetStore) AddToMatchupSet(ctx context.Context, contender1 string) error {
	otherContenders := &Contenders{}
	if err := s.db.Scan(ctx, otherContenders); err != nil {
		return errors.Wrap(err, "failed to retrieve other contenders to populate Matchup Set")
	}

	for _, contender2 := range otherContenders {
		c1, c2 := OrderMatchup(contender1, contender2)
		newMatchupEntry := matchupSetEntry{
			Contender1: c1,
			Contender2: c2,
		}

		matchupSet := &MatchupSet{
			ID:        masterKey,
			entry:     newMatchupEntry,
			tableName: masterMatchupTableName,
		}
		if err := s.db.Update(ctx, matchupSet); err != nil {
			return errors.Wrapf(err, "failed to update the matchup set for ID: %s", masterKey)
		}

	}
	return nil
}

// RemoveFromMatchupSet given a uid corresponding to the session, and two contenders, adds them to the set
// of matchups that uid has seen
func (s *MasterMatchupSetStore) RemoveFromMatchupSet(ctx context.Context, uid, contender1, contender2 string) error {
	contender1, contender2 = OrderMatchup(contender1, contender2)
	newMatchupEntry := matchupSetEntry{
		Contender1: contender1,
		Contender2: contender2,
		remove:     true,
	}

	matchupSet := &MatchupSet{
		ID:        uid,
		entry:     newMatchupEntry,
		tableName: masterMatchupTableName,
	}
	if err := s.db.Update(ctx, matchupSet); err != nil {
		return errors.Wrapf(err, "failed to update the matchup set for ID: %s", uid)
	}
	return nil
}
