package dynamostore

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

// TableConfig allows us to set configuration details
// for the dynamo table from the app
type TableConfig struct {
	TableName        string
	ReadCapacity     int
	WriteCapacity    int
	TTLEnabled       bool
	TTLAttributeName string
}

// Flags returns a slice of the configuration options for the contender table
// it assumes that the prefix passed in is separated by dashes, if at all
func (c *TableConfig) Flags(prefix, defaultTableName string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        cliFlagName(prefix, "table-name"),
			EnvVar:      envVarName(prefix, "TABLE_NAME"),
			Value:       defaultTableName,
			Destination: &c.TableName,
		},
		cli.IntFlag{
			Name:        cliFlagName(prefix, "table-read-capacity"),
			EnvVar:      envVarName(prefix, "TABLE_READ_CAPACITY"),
			Value:       5,
			Destination: &c.ReadCapacity,
		},
		cli.IntFlag{
			Name:        cliFlagName(prefix, "table-write-capacity"),
			EnvVar:      envVarName(prefix, "TABLE_WRITE_CAPACITY"),
			Value:       5,
			Destination: &c.WriteCapacity,
		},
		cli.BoolFlag{
			Name:        cliFlagName(prefix, "table-ttl-enabled"),
			EnvVar:      envVarName(prefix, "TABLE_TTL_ENABLED"),
			Destination: &c.TTLEnabled,
		},
		cli.StringFlag{
			Name:        cliFlagName(prefix, "table-ttl-attribute-name"),
			EnvVar:      envVarName(prefix, "TABLE_TTL_ATTRIBUTE_NAME"),
			Destination: &c.TTLAttributeName,
		},
	}
}

func envVarName(prefix, name string) string {
	return strings.Replace("-", "_", strings.ToUpper(cliFlagName(prefix, name)), -1)
}

func cliFlagName(prefix, name string) string {
	return fmt.Sprintf("%s-%s", prefix, name)
}
