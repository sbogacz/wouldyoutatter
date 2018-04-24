package contender

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
	"github.com/urfave/cli"
)

// LeaderboardEntry is the model for the head-to-head records
// between contenders
type LeaderboardEntry struct {
	Contender   string
	Score       int
	Wins        int
	entrantLost bool
}

// Leaderboard is a collection of LeaderboardEntrys
type Leaderboard []LeaderboardEntry

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
func (s *LeaderboardStore) Set(ctx context.Context, m *LeaderboardEntry) error {
	return errors.Wrap(s.db.Set(ctx, m), "failed to save leaderboard")
}

// Get retrieves the entire leaderboard
func (s *LeaderboardStore) Get(ctx context.Context) (Leaderboard, error) {
	leaderboardEntries := []LeaderboardEntry{}
	leaderboard := Leaderboard(leaderboardEntries)
	if err := s.db.Scan(ctx, &leaderboard); err != nil {
		return nil, errors.Wrap(err, "failed to retrieve other contenders to populate Matchup Set")
	}
	return leaderboard, nil
}

// GetEntry lets you retrieve a leaderboard entry by the contender name
func (s *LeaderboardStore) GetEntry(ctx context.Context, contender string) (*LeaderboardEntry, error) {
	m := &LeaderboardEntry{Contender: contender}
	item, err := s.db.Get(ctx, m)
	if err != nil {
		if dynamostore.NotFoundError(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to retrieve leaderboard")

	}
	ret := item.(*LeaderboardEntry)
	return ret, nil
}

// Delete lets you delete a leaderboard entry by name
func (s *LeaderboardStore) Delete(ctx context.Context, contender string) error {
	m := &LeaderboardEntry{Contender: contender}

	return errors.Wrap(s.db.Delete(ctx, m), "failed to delete leaderboard")
}

// UpdateWinningEntry lets you declarea a leaderboard entry a winner by name
func (s *LeaderboardStore) UpdateWinningEntry(ctx context.Context, name string, oldScore int) error {
	winner := &LeaderboardEntry{Contender: name, Score: oldScore}

	return errors.Wrapf(s.db.Update(ctx, winner), "failed to declare leaderboard %s the winner", name)
}

// UpdateLosingEntry lets you declarea a leaderboard entry a loser by name
func (s *LeaderboardStore) UpdateLosingEntry(ctx context.Context, name string, oldScore int) error {
	loser := &LeaderboardEntry{Contender: name, Score: oldScore, entrantLost: true}

	return errors.Wrapf(s.db.Update(ctx, loser), "failed to declare leaderboard %s the loser", name)
}

// LeaderboardTableConfig allows us to set configuration details
// for the dynamo table from the app
type LeaderboardTableConfig struct {
	TableName     string
	ReadCapacity  int
	WriteCapacity int
}

// Flags returns a slice of the configuration options for the leaderboard table
func (c *LeaderboardTableConfig) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "leaderboard-table-name",
			EnvVar:      "LEADERBOARD_TABLE_NAME",
			Value:       "Leaderboard",
			Destination: &c.TableName,
		},
		cli.IntFlag{
			Name:        "leaderboard-table-read-capacity",
			EnvVar:      "LEADERBOARD_TABLE_READ_CAPACITY",
			Value:       5,
			Destination: &c.ReadCapacity,
		},
		cli.IntFlag{
			Name:        "leaderboard-table-write-capacity",
			EnvVar:      "LEADERBOARD_TABLE_WRITE_CAPACITY",
			Value:       5,
			Destination: &c.WriteCapacity,
		},
	}
}
