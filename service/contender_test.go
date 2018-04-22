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
