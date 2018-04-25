package service

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	// DefaultPort for the service
	DefaultPort = 8080
	// DefaultMasterKey for the service
	DefaultMasterKey = "th3M0stm3tAlTh1ng1Hav3ev3rh3ard"
	// DefaultLogLevel for the service
	DefaultLogLevel = "INFO"

	// DefaultContenderTableName is what it sounds like
	DefaultContenderTableName = "Contenders"
	// DefaultMatchupTableName is what it sounds like
	DefaultMatchupTableName = "Matchups"
	// DefaultUserMatchupsTableName is what it sounds like
	DefaultUserMatchupsTableName = "User-Past-Matchups"
	// DefaultMasterMatchupsTableName is what it sounds like
	DefaultMasterMatchupsTableName = "Possible-Matchups"
	// DefaultTokenTableName is what it sounds like
	DefaultTokenTableName = "Tokens"
)

var (
	awsConfig *aws.Config
)

// Config holds the service variables we want to
// to configure from the cli/env
type Config struct {
	Port           int
	AWSAccessKeyID string
	AWSSecretKey   string
	AWSRegion      string
	MasterKey      string
	LogLevel       string
	// Table Configs
	ContenderTableConfig      *dynamostore.TableConfig
	MatchupTableConfig        *dynamostore.TableConfig
	UserMatchupsTableConfig   *dynamostore.TableConfig
	MasterMatchupsTableConfig *dynamostore.TableConfig
	TokenTableConfig          *dynamostore.TableConfig
}

// Flags r	eturns the slice of cli.Flags that we have
// available
func (c *Config) Flags() []cli.Flag {
	ret := []cli.Flag{
		cli.IntFlag{
			Name:        "port, p",
			Usage:       "the port you'd like to run the service on",
			Destination: &c.Port,
			Value:       DefaultPort,
		},
		cli.StringFlag{
			Name:        "aws-access-key-id",
			EnvVar:      "AWS_ACCESS_KEY_ID",
			Usage:       "the AWS access key to sign requests with",
			Destination: &c.AWSAccessKeyID,
		},
		cli.StringFlag{
			Name:        "aws-secret-key",
			EnvVar:      "AWS_SECRET_ACCESS_KEY",
			Usage:       "the AWS secret key to sign requests with",
			Destination: &c.AWSSecretKey,
		},
		cli.StringFlag{
			Name:        "aws-region",
			EnvVar:      "AWS_REGION",
			Usage:       "the AWS region to connect to",
			Destination: &c.AWSRegion,
		},
		cli.StringFlag{
			Name:        "master-key",
			EnvVar:      "MASTER_KEY",
			Usage:       "the master key that will be required in the header for Creates and Deletes on contenders",
			Destination: &c.MasterKey,
			Value:       DefaultMasterKey,
		},
		cli.StringFlag{
			Name:        "log-level",
			EnvVar:      "LOG_LEVEL",
			Usage:       "the log level that wouldyoutatter should log to stdout at, defaults to INFO",
			Destination: &c.LogLevel,
			Value:       DefaultLogLevel,
		},
	}
	// initialize configs
	c.ContenderTableConfig = &dynamostore.TableConfig{}
	c.MatchupTableConfig = &dynamostore.TableConfig{}
	c.UserMatchupsTableConfig = &dynamostore.TableConfig{}
	c.MasterMatchupsTableConfig = &dynamostore.TableConfig{}
	c.TokenTableConfig = &dynamostore.TableConfig{}

	ret = append(ret, c.ContenderTableConfig.Flags("contender", DefaultContenderTableName)...)
	ret = append(ret, c.MatchupTableConfig.Flags("matchup", DefaultMatchupTableName)...)
	ret = append(ret, c.UserMatchupsTableConfig.Flags("user-matchups", DefaultUserMatchupsTableName)...)
	ret = append(ret, c.MasterMatchupsTableConfig.Flags("master-matchups", DefaultMasterMatchupsTableName)...)
	ret = append(ret, c.TokenTableConfig.Flags("token", DefaultTokenTableName)...)
	return ret
}

// Retrieve allows the Config to be used as an aws.Provider
func (c *Config) Retrieve() (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     c.AWSAccessKeyID,
		SecretAccessKey: c.AWSSecretKey,
		CanExpire:       false,
	}, nil
}

func (c *Config) logLevelToLogrus() log.Level {
	switch c.LogLevel {
	case "DEBUG":
		return log.DebugLevel
	case "WARN":
		return log.WarnLevel
	case "ERROR":
		return log.ErrorLevel
	default:
	}
	return log.InfoLevel
}
