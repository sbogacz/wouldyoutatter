package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

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

		leaderboard := contender.Leaderboard{}
		d := json.NewDecoder(resp.Body)
		err = d.Decode(&leaderboard)
		require.NoError(t, err)

		assert.Equal(t, len(clientSideLeaderboard), len(leaderboard))
		// check that the leaderboard's scores match ours

		fmt.Printf("%+v\n", clientSideLeaderboard)
		fmt.Printf("%+v\n", clientSideWins)
		time.Sleep(time.Minute)
		for _, entry := range leaderboard {
			fmt.Println("entry " + entry.Contender)
			assert.Equal(t, clientSideLeaderboard[entry.Contender], entry.Score)
			assert.Equal(t, clientSideWins[entry.Contender], entry.Wins)
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
