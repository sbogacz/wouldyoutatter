package service_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/phayes/freeport"
	"github.com/sbogacz/wouldyoutatter/service"
	log "github.com/sirupsen/logrus"
)

var (
	runAgainstLocalDynamo = flag.Bool("local-dynamo", false, "run tests against a local dynamo instance")
	s                     *service.Service
	baseAddress           string
	contenderAddress      string
	matchupAddress        string
	leaderboardAddress    string
)

func TestMain(m *testing.M) {
	flag.Parse()

	// set up
	openPort, err := freeport.GetFreePort()
	if err != nil {
		log.Fatalf("failed to get free port for tests: %v", err)
	}
	config := service.Config{}

	// configure default table names
	for _, f := range config.Flags() {
		f.Apply(flag.CommandLine)
	}
	// override options for the test
	config.Port = openPort
	config.LogLevel = "INFO"

	if *runAgainstLocalDynamo {
		config.AWSRegion = "local"
	}

	if err := setupService(config); err != nil {
		log.Fatalf("failed to setup for tests: %v", err)
		return
	}
	baseAddress = fmt.Sprintf("http://127.0.0.1:%d", openPort)
	contenderAddress = fmt.Sprintf("%s/contenders", baseAddress)
	matchupAddress = fmt.Sprintf("%s/matchups", baseAddress)
	leaderboardAddress = fmt.Sprintf("%s/leaderboard", baseAddress)

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

	tables := []string{
		service.DefaultContenderTableName,
		service.DefaultMasterMatchupsTableName,
		service.DefaultUserMatchupsTableName,
		service.DefaultTokenTableName,
		//service.DefaultMatchupTableName,
	}

	for _, table := range tables {
		input := &dynamodb.DeleteTableInput{
			TableName: aws.String(table),
		}
		log.WithField("tablename", table).Info("deleting table after tests")
		req := svc.DeleteTableRequest(input)
		_, err = req.Send()
		if err != nil {
			return err
		}
	}
	return nil
}
