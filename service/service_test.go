package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/phayes/freeport"
	"github.com/sbogacz/wouldyoutatter/contender"
	"github.com/sbogacz/wouldyoutatter/service"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	runIntegrationTests = true
	s                   *service.Service
	baseAddress         string
	contenderAddress    string
)

func TestMain(m *testing.M) {
	if err := setupService(); err != nil {
		log.Fatalf("failed to setup for tests: %v", err)
		return
	}
	go s.Start()
	status := m.Run()
	s.Stop()
	os.Exit(status)
}

func TestSimpleContenderCRUD(t *testing.T) {
	origContender := contender.Contender{
		Name:        "banana",
		Description: "an apple",
	}
	t.Run("create contender", func(t *testing.T) {
		b, err := json.Marshal(&origContender)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Post(contenderAddress, "application/json", bytes.NewBuffer(b))
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})
	t.Run("get contender", func(t *testing.T) {
		resp, err := http.DefaultClient.Get(fmt.Sprintf("%s/%s", contenderAddress, origContender.Name))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		require.NotNil(t, resp.Body)
		d := json.NewDecoder(resp.Body)

		storedContender := contender.Contender{}
		err = d.Decode(&storedContender)
		require.NoError(t, err)
		assert.Equal(t, storedContender.Description, origContender.Description)
	})
	t.Run("delete contender", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", contenderAddress, origContender.Name), nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp, err = http.DefaultClient.Get(fmt.Sprintf("%s/%s", contenderAddress, origContender.Name))
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func setupService() error {
	openPort, err := freeport.GetFreePort()
	if err != nil {
		return err
	}
	config := service.Config{Port: openPort}
	s = service.New(config, contender.NewLocalStore())
	baseAddress = fmt.Sprintf("http://127.0.0.1:%d", openPort)
	contenderAddress = fmt.Sprintf("%s/contenders", baseAddress)
	return nil
}
