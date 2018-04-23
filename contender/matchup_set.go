package contender

import (
	"context"
	"fmt"
	"strings"

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
	ID    string
	Set   []MatchupSetEntry
	entry MatchupSetEntry
}

// MatchupSetEntry holds a possible matchup combination
type MatchupSetEntry struct {
	Contender1 string
	Contender2 string
	VoteURL    string // we don't record this in the DB, but we use it in the API
	remove     bool
}

func newMatchupSetEntry(c1, c2 string) MatchupSetEntry {
	contender1, contender2 := OrderMatchup(c1, c2)
	return MatchupSetEntry{
		Contender1: contender1,
		Contender2: contender2,
	}
}

func matchupEntryfromString(s string) MatchupSetEntry {
	contenders := strings.Split(s, "§")
	contender1, contender2 := OrderMatchup(contenders[0], contenders[1])
	return MatchupSetEntry{
		Contender1: contender1,
		Contender2: contender2,
	}
}

func (m MatchupSetEntry) String() string {
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

// Add given a uid corresponding to the session, and two contenders, adds them to the set
// of matchups that uid has seen
func (s *MatchupSetStore) Add(ctx context.Context, uid, contender1, contender2 string) error {
	newEntry := newMatchupSetEntry(contender1, contender2)

	matchupSet := &MatchupSet{
		ID:    uid,
		entry: newEntry,
	}
	if err := s.db.Update(ctx, matchupSet); err != nil {
		return errors.Wrapf(err, "failed to update the matchup set for ID: %s", uid)
	}
	return nil
}

// Remove given a uid corresponding to the session, and two contenders, adds them to the set
// of matchups that uid has seen
func (s *MatchupSetStore) Remove(ctx context.Context, uid, contender1, contender2 string) error {
	contender1, contender2 = OrderMatchup(contender1, contender2)
	newMatchupEntry := MatchupSetEntry{
		Contender1: contender1,
		Contender2: contender2,
		remove:     true,
	}

	matchupSet := &MatchupSet{
		ID:    uid,
		entry: newMatchupEntry,
	}
	if err := s.db.Update(ctx, matchupSet); err != nil {
		return errors.Wrapf(err, "failed to update the matchup set for ID: %s", uid)
	}
	return nil
}

// Get lets you retrieve a contender by name
func (s *MatchupSetStore) Get(ctx context.Context, uid string) (*MatchupSet, error) {
	m := &MatchupSet{ID: uid}
	item, err := s.db.Get(ctx, m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve matchup set")

	}
	ret := item.(*MatchupSet)
	return ret, nil
}

// Delete is used to restart a matchup set when it is no longer relevant
func (s *MatchupSetStore) Delete(ctx context.Context, uid string) error {
	if err := s.db.Delete(ctx, &MatchupSet{ID: uid}); err != nil {
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

// Add given a uid corresponding to the session, and two contenders, adds them to the set
// of matchups that uid has seen
func (s *MasterMatchupSetStore) Add(ctx context.Context, contender1 string, otherContenders *Contenders) error {
	for _, contender2 := range *otherContenders {
		// don't create dupes
		if contender1 == contender2.Name {
			continue
		}
		c1, c2 := OrderMatchup(contender1, contender2.Name)
		newMatchupEntry := MatchupSetEntry{
			Contender1: c1,
			Contender2: c2,
		}

		matchupSet := &MatchupSet{
			ID:    masterKey,
			entry: newMatchupEntry,
		}
		if err := s.db.Update(ctx, matchupSet); err != nil {
			return errors.Wrapf(err, "failed to update the matchup set for ID: %s", masterKey)
		}

	}
	return nil
}

// Remove given a uid corresponding to the session, and two contenders, adds them to the set
// of matchups that uid has seen
func (s *MasterMatchupSetStore) Remove(ctx context.Context, contender1 string, otherContenders *Contenders) error {

	for _, contender2 := range *otherContenders {
		// didn't create dupes
		if contender1 == contender2.Name {
			continue
		}
		c1, c2 := OrderMatchup(contender1, contender2.Name)
		newMatchupEntry := MatchupSetEntry{
			Contender1: c1,
			Contender2: c2,
			remove:     true,
		}

		matchupSet := &MatchupSet{
			ID:    masterKey,
			entry: newMatchupEntry,
		}
		if err := s.db.Update(ctx, matchupSet); err != nil {
			return errors.Wrapf(err, "failed to update the matchup set for ID: %s", masterKey)
		}

	}
	return nil
}

// Get lets you retrieve a contender by name
func (s *MasterMatchupSetStore) Get(ctx context.Context) (*MatchupSet, error) {
	m := &MatchupSet{ID: masterKey}
	item, err := s.db.Get(ctx, m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve matchup set")

	}
	ret := item.(*MatchupSet)
	return ret, nil
}
