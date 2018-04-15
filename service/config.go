package service

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var (
	awsConfig *aws.Config
)

// Config holds the service variables we want to
// to configure from the cli/env
type Config struct {
	Port         int
	AWSSecretKey string
	AWSAccessKey string
	AWSRegion    string
	AWSEndpoint  string
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
			Name:        "aws-secret-key",
			EnvVar:      "AWS_SECRET_ACCESS_KEY",
			Usage:       "the AWS secret key to sign requests with",
			Destination: &c.AWSSecretKey,
		},
	}
}

func (c *Config) GetAWSConfig() error {
	_, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return errors.Wrap(err, "failed to load default AWS config")
	}
	return nil
}
