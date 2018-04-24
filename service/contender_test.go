package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/sbogacz/wouldyoutatter/contender"
	"github.com/sbogacz/wouldyoutatter/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleContenderCRUD(t *testing.T) {
	origContender := contender.Contender{
		Name:        "banana",
		Description: "an apple",
		SVG:         []byte("pretend this is an svg"),
	}
	t.Run("create contender", func(t *testing.T) {
		b, err := json.Marshal(&origContender)
		require.NoError(t, err)

		// first go should fail since it's unauthorized
		req, err := http.NewRequest("POST", contenderAddress, bytes.NewBuffer(b))
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// if we add the master key header it should succeed
		req, err = http.NewRequest("POST", contenderAddress, bytes.NewBuffer(b))
		require.NoError(t, err)
		req.Header.Set("X-Tatter-Master", service.DefaultMasterKey)
		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

	})
	t.Run("get contender", func(t *testing.T) {
		resp, err := http.DefaultClient.Get(fmt.Sprintf("%s/%s", contenderAddress, origContender.Name))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		require.NotNil(t, resp.Body)
		d := json.NewDecoder(resp.Body)

		storedContender := contender.Contender{}
		err = d.Decode(&storedContender)
		require.NoError(t, err)
		require.Equal(t, storedContender.Description, origContender.Description)
	})
	t.Run("delete contender", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", contenderAddress, origContender.Name), nil)
		require.NoError(t, err)

		// should fail first time, unauthorized
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// try again with right header
		req, err = http.NewRequest("DELETE", fmt.Sprintf("%s/%s", contenderAddress, origContender.Name), nil)
		require.NoError(t, err)
		req.Header.Set("X-Tatter-Master", service.DefaultMasterKey)

		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// try a get again and make sure we 404
		resp, err = http.DefaultClient.Get(fmt.Sprintf("%s/%s", contenderAddress, origContender.Name))
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestAddingSeveralContendersCreatesPossibleMatchups(t *testing.T) {
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

	t.Run("as we ask for matchups, we should be able to see 6 different ones before looping", func(t *testing.T) {
		var cookie *http.Cookie
		previousMatchups := []string{}
		var sawRepeat bool

		for {
			if sawRepeat {
				if len(previousMatchups) != 6 {
					fmt.Printf("saw a repeat without seeing every combination %d\n", len(previousMatchups))
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

			if stringInSlice(matchup.String(), previousMatchups) {
				sawRepeat = true
				continue
			}
			previousMatchups = append(previousMatchups, matchup.String())
		}
		require.True(t, true)

	})
}

func stringInSlice(s string, arr []string) bool {
	for _, str := range arr {
		if s == str {
			return true
		}
	}
	return false
}
