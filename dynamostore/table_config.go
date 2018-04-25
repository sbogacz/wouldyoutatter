package dynamostore

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

// TableConfig allows us to set configuration details
// for the dynamo table from the app
type TableConfig struct {
	TableName     string
	ReadCapacity  int64
	WriteCapacity int64
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
		cli.Int64Flag{
			Name:        cliFlagName(prefix, "table-read-capacity"),
			EnvVar:      envVarName(prefix, "TABLE_READ_CAPACITY"),
			Value:       5,
			Destination: &c.ReadCapacity,
		},
		cli.Int64Flag{
			Name:        cliFlagName(prefix, "table-write-capacity"),
			EnvVar:      envVarName(prefix, "TABLE_WRITE_CAPACITY"),
			Value:       5,
			Destination: &c.WriteCapacity,
		},
	}
}

func envVarName(prefix, name string) string {
	return strings.Replace("-", "_", strings.ToUpper(cliFlagName(prefix, name)), -1)
}

func cliFlagName(prefix, name string) string {
	return fmt.Sprintf("%s-%s", prefix, name)
}
