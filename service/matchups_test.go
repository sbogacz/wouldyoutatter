package service_test

import (
	"bytes"
	"container/heap"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/sbogacz/wouldyoutatter/contender"
	"github.com/sbogacz/wouldyoutatter/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVotingAndLeaderboard(t *testing.T) {
	things := []string{"banana", "apple", "window", "troll"}
	contenders := []contender.Contender{}
	for _, thing := range things {
		contenders = append(contenders, contender.Contender{
			Name:        fmt.Sprintf("%s", thing),
			Description: fmt.Sprintf("a %s", thing),
			SVG:         []byte(fmt.Sprintf("pretend this is an svg of %s", thing)),
		})

	}

	t.Run("create all of the contenders", func(t *testing.T) {
		for _, contender := range contenders {
			b, err := json.Marshal(&contender)
			require.NoError(t, err)

			// first go should fail since it's unauthorized
			req, err := http.NewRequest("POST", contenderAddress, bytes.NewBuffer(b))
			require.NoError(t, err)

			req.Header.Set("X-Tatter-Master", service.DefaultMasterKey)
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusCreated, resp.StatusCode)
		}
	})
	// we'll popuulate this set in the next step to be used in the voting step
	matchups := []contender.MatchupSetEntry{}

	t.Run("as we ask for matchups, we should be able to see 6 different ones before looping", func(t *testing.T) {
		var cookie *http.Cookie
		var sawRepeat bool

		for {
			if sawRepeat {
				if len(matchups) != 6 {
					//time.Sleep(time.Minute)
					require.True(t, false)
				}
				break
			}
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/random", matchupAddress), nil)
			require.NoError(t, err)
			// if have cookie, set
			if cookie != nil {
				req.AddCookie(cookie)
			}
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			// if we didn't have a cookie, set it from the response
			if cookie == nil {
				for _, c := range resp.Cookies() {
					if c.Name == service.CookieKey {
						cookie = c
					}
				}
			}

			matchup := &contender.MatchupSetEntry{}
			d := json.NewDecoder(resp.Body)
			err = d.Decode(matchup)
			require.NoError(t, err)

			resp.Body.Close()

			if matchupInSlice(matchup, matchups) {
				sawRepeat = true
				continue
			}
			matchups = append(matchups, *matchup)
		}
		require.True(t, true)
	})

	clientSideLeaderboard := make(map[string]int, len(contenders))
	clientSideWins := make(map[string]int, len(contenders))
	t.Run("loop through the matchup list and use the URL to vote", func(t *testing.T) {
		for i, matchup := range matchups {
			// for the first three, vote for the first contender
			winner := matchup.Contender2
			loser := matchup.Contender1
			if i < 3 {
				winner = matchup.Contender1
				loser = matchup.Contender2
			}
			// update our client side leaderboard to use for later verification
			clientSideLeaderboard[winner] = clientSideLeaderboard[winner] + 1
			clientSideWins[winner] = clientSideWins[winner] + 1
			clientSideLeaderboard[loser] = clientSideLeaderboard[loser] - 1

			payload := votePayload{Winner: winner}
			b, err := json.Marshal(&payload)
			require.NoError(t, err)

			u := fmt.Sprintf("%s%s", baseAddress, matchup.VoteURL)
			resp, err := http.DefaultClient.Post(u, "application/json", bytes.NewReader(b))

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("check the leaderboard", func(t *testing.T) {
		resp, err := http.DefaultClient.Get(leaderboardAddress)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		dbLeaderboard := contender.Contenders{}
		d := json.NewDecoder(resp.Body)
		err = d.Decode(&dbLeaderboard)
		require.NoError(t, err)

		assert.Equal(t, len(clientSideLeaderboard), len(dbLeaderboard))
		// check that the leaderboard's scores match ours

		for _, c := range dbLeaderboard {
			assert.Equal(t, clientSideLeaderboard[c.Name], c.Score)
			assert.Equal(t, clientSideWins[c.Name], c.Wins)
		}
	})

	t.Run("check the leaderboard top 3", func(t *testing.T) {
		// first create our local top 3
		localLeaderboard := &leaderboard{}
		heap.Init(localLeaderboard)
		for k, v := range clientSideLeaderboard {
			heap.Push(localLeaderboard, entry{name: k, score: v})
		}

		resp, err := http.DefaultClient.Get(fmt.Sprintf("%s?limit=3", leaderboardAddress))
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		top3 := contender.Contenders{}
		d := json.NewDecoder(resp.Body)
		err = d.Decode(&top3)
		require.NoError(t, err)

		assert.Equal(t, 3, len(top3))

		// check that the leaderboard's top 3scores match ours
		for _, c := range top3 {
			localEntry := heap.Pop(localLeaderboard).(entry)
			assert.Equal(t, localEntry.name, c.Name)
			assert.Equal(t, localEntry.score, c.Score)
		}
	})
}

type votePayload struct {
	Winner string
}

func matchupInSlice(m *contender.MatchupSetEntry, arr []contender.MatchupSetEntry) bool {
	for _, matchup := range arr {
		if m.String() == matchup.String() {
			return true
		}
	}
	return false
}

type entry struct {
	score int
	name  string
}

type leaderboard []entry

func (l leaderboard) Len() int           { return len(l) }
func (l leaderboard) Less(i, j int) bool { return l[i].score > l[j].score }
func (l leaderboard) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func (l *leaderboard) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the shice's hength,
	// not just its contents.
	*l = append(*l, x.(entry))
}

func (l *leaderboard) Pop() interface{} {
	old := *l
	n := len(old)
	x := old[n-1]
	*l = old[0 : n-1]
	return x
}
