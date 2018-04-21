package service_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/phayes/freeport"
	"github.com/sbogacz/wouldyoutatter/contender"
	"github.com/sbogacz/wouldyoutatter/service"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	runAgainstLocalDynamo = flag.Bool("local-dynamo", false, "run tests against a local dynamo instance")
	s                     *service.Service
	baseAddress           string
	contenderAddress      string
)

func TestMain(m *testing.M) {
	flag.Parse()

	// set up
	openPort, err := freeport.GetFreePort()
	if err != nil {
		log.Fatalf("failed to get free port for tests: %v", err)
	}
	config := service.Config{
		Port:      openPort,
		LogLevel:  "INFO",
		MasterKey: service.DefaultMasterKey,
	}

	if *runAgainstLocalDynamo {
		config.AWSRegion = "local"
	}

	if err := setupService(config); err != nil {
		log.Fatalf("failed to setup for tests: %v", err)
		return
	}
	baseAddress = fmt.Sprintf("http://127.0.0.1:%d", openPort)
	contenderAddress = fmt.Sprintf("%s/contenders", baseAddress)

	go s.Start()
	status := m.Run()
	s.Stop()

	// tear down
	if *runAgainstLocalDynamo {
		config.AWSRegion = "local"
		if err := teardownTables(config); err != nil {
			log.Fatalf("failed to tear down tables for tests: %v", err)
		}
	}
	os.Exit(status)
}

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

func setupService(config service.Config) error {
	var err error

	s, err = service.New(config)
	if err != nil {
		return err
	}
	return nil
}

func teardownTables(config service.Config) error {
	cfg, err := config.AWSConfig()
	if err != nil {
		return err
	}
	svc := dynamodb.New(cfg)
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String("Contenders"),
	}

	req := svc.DeleteTableRequest(input)
	_, err = req.Send()
	if err != nil {
		return err
	}
	return nil
}
