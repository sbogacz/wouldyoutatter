package service

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/urfave/cli"
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
}

// Flags r	eturns the slice of cli.Flags that we have
// available
func (c *Config) Flags() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:        "port, p",
			Usage:       "the port you'd like to run the service on",
			Destination: &c.Port,
			Value:       8080,
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
	}
}

// AWSConfig returns an aws Config based on the env vars/flags
// provided
func (c *Config) AWSConfig() (aws.Config, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return aws.Config{}, err
	}
	cfg.Region = c.AWSRegion
	cfg.Credentials = c
	return cfg, nil
}

// Retrieve allows the Config to be used as an aws.Provider
func (c *Config) Retrieve() (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     c.AWSAccessKeyID,
		SecretAccessKey: c.AWSSecretKey,
		CanExpire:       false,
	}, nil
}
