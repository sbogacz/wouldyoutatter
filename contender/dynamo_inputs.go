package contender

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
)

var _ dynamostore.Item = (*Contender)(nil)

const (
	leaderboardScoreIndex = "LeaderboardScore"
)

// Key returns the Contenders name, and implements the dynamostore Item interface
func (c Contender) Key() string {
	return c.Name
}

// Marshal encodes the values of a contender into the map format
// that dynamo expects
func (c Contender) Marshal() map[string]dynamodb.AttributeValue {
	return map[string]dynamodb.AttributeValue{
		"Name":        stringToAttributeValue(c.Name),
		"Description": stringToAttributeValue(c.Description),
		"SVG":         bytesToAttributeValue(c.SVG),
		"Wins":        intToAttributeValue(c.Wins),
		"Losses":      intToAttributeValue(c.Losses),
		"Score":       intToAttributeValue(c.Score),
		"Leaderboard": stringToAttributeValue("topscore"), // placeholder
	}
}

// Unmarshal tries to decode a Contender from a dynamo response
func (c *Contender) Unmarshal(aMap map[string]dynamodb.AttributeValue) error {
	if len(aMap) == 0 {
		return errors.New(dynamodb.ErrCodeResourceNotFoundException)
	}

	wins, err := getInt(aMap["Wins"])
	if err != nil {
		return errors.Wrap(err, "failed to read Wins attribute")
	}
	losses, err := getInt(aMap["Losses"])
	if err != nil {
		return errors.Wrap(err, "failed to read Losses attribute")
	}
	score, err := getInt(aMap["Score"])
	if err != nil {
		return errors.Wrap(err, "failed to read Score attribute")
	}
	newContender := &Contender{
		Name:        getString(aMap["Name"]),
		Description: getString(aMap["Description"]),
		SVG:         getBytes(aMap["SVG"]),
		Wins:        wins,
		Losses:      losses,
		Score:       score,
	}
	*c = *newContender
	return nil
}

// CreateTableInput generates the dynamo input to create the contenders table
func (c *Contender) CreateTableInput(tc *dynamostore.TableConfig) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Name"),
				AttributeType: dynamodb.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("Leaderboard"),
				AttributeType: dynamodb.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("Score"),
				AttributeType: dynamodb.ScalarAttributeTypeN,
			},
		},
		KeySchema: []dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Name"),
				KeyType:       dynamodb.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(tc.ReadCapacity),
			WriteCapacityUnits: aws.Int64(tc.WriteCapacity),
		},
		GlobalSecondaryIndexes: []dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String(leaderboardScoreIndex),
				KeySchema: []dynamodb.KeySchemaElement{
					{
						// placeholder to allow us to sort our results by
						AttributeName: aws.String("Leaderboard"),
						KeyType:       dynamodb.KeyTypeHash,
					},
					{
						AttributeName: aws.String("Score"),
						KeyType:       dynamodb.KeyTypeRange,
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(tc.ReadCapacity),
					WriteCapacityUnits: aws.Int64(tc.WriteCapacity),
				},
				Projection: &dynamodb.Projection{
					ProjectionType: dynamodb.ProjectionTypeAll,
				},
			},
		},

		TableName: aws.String(tc.TableName),
	}
}

// DescribeTableInput generates the query we need to describe the contender table
func (c *Contender) DescribeTableInput(tableName string) *dynamodb.DescribeTableInput {
	return &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
}

// TableOptions is a no-op for the contender table (for now)
func (c *Contender) TableOptions(tableName string) []dynamostore.TableOption {
	return nil
}

// GetItemInput generates the dynamodb.GetItemInput for the given contender
func (c *Contender) GetItemInput(tableName string) *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       map[string]dynamodb.AttributeValue{"Name": {S: aws.String(c.Name)}},
	}
}

// PutItemInput generates the dynamodb.PutItemInput for the given contender
func (c *Contender) PutItemInput(tableName string) *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      c.Marshal(),
	}
}

// DeleteItemInput generates the dynamodb.DeleteItemInput for the given contender
func (c *Contender) DeleteItemInput(tableName string) *dynamodb.DeleteItemInput {
	return &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       map[string]dynamodb.AttributeValue{"Name": {S: aws.String(c.Name)}},
	}
}

// UpdateItemInput generates the dynamodb.UpdateItemInput for the given contender
func (c *Contender) UpdateItemInput(tableName string) *dynamodb.UpdateItemInput {
	if c.isLoser {
		return lossInput(c.Name, tableName)
	}
	return winInput(c.Name, tableName)
}

func winInput(name, tableName string) *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       map[string]dynamodb.AttributeValue{"Name": {S: aws.String(name)}},
		UpdateExpression:          aws.String("ADD Wins :w, Score :w"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{":w": {N: aws.String("1")}},
	}
}

func lossInput(name, tableName string) *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName:        aws.String(tableName),
		Key:              map[string]dynamodb.AttributeValue{"Name": {S: aws.String(name)}},
		UpdateExpression: aws.String("ADD Losses :l, Score :ls"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":l":  {N: aws.String("1")},
			":ls": {N: aws.String("-1")},
		},
	}
}

// ScanInput produces a dynamodb ScanInput object
func (c *Contenders) ScanInput(tableName string) *dynamodb.ScanInput {
	return &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
}

// QueryInput producest a dynamodb QueryInput object looking for the
// top N contenders
func (c *Contenders) QueryInput(tableName string, limit int) *dynamodb.QueryInput {
	return &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(leaderboardScoreIndex),
		KeyConditionExpression:    aws.String("Leaderboard = :val"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{":val": {S: aws.String("topscore")}},
		Limit:            aws.Int64(int64(limit)),
		ScanIndexForward: aws.Bool(false),
	}
}

// Unmarshal allows results to be unmarshalled directly into the struct
func (c *Contenders) Unmarshal(maps []map[string]dynamodb.AttributeValue) error {
	cs := make([]*Contender, len(maps))
	for i := range cs {
		cs[i] = &Contender{}
		if err := cs[i].Unmarshal(maps[i]); err != nil {
			return errors.Wrap(err, "failed to unmarshal Contenders")
		}
	}
	contenders := make([]Contender, len(cs))
	for i := range cs {
		contenders[i] = *cs[i]
	}
	*c = contenders
	return nil

}

func stringToAttributeValue(s string) dynamodb.AttributeValue {
	return dynamodb.AttributeValue{S: aws.String(s)}
}

func intToAttributeValue(n int) dynamodb.AttributeValue {
	return dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%d", n))}
}

func int64ToAttributeValue(n int64) dynamodb.AttributeValue {
	return dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%d", n))}
}

func bytesToAttributeValue(b []byte) dynamodb.AttributeValue {
	return dynamodb.AttributeValue{B: b}
}

func getString(a dynamodb.AttributeValue) string {
	if a.S == nil {
		return ""
	}
	return *a.S
}

func getInt(a dynamodb.AttributeValue) (int, error) {
	if a.N == nil {
		return 0, nil
	}
	return strconv.Atoi(*a.N)
}

func getBytes(a dynamodb.AttributeValue) []byte {
	return a.B
}
