package contender

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
)

var _ dynamostore.Item = (*LeaderboardEntry)(nil)

// Key returns the LeaderboardEntrys name, and implements the dynamostore Item interface
func (l LeaderboardEntry) Key() string {
	return l.Contender
}

// Marshal encodes the values of a contender into the map format
// that dynamo expects
func (l LeaderboardEntry) Marshal() map[string]dynamodb.AttributeValue {
	return map[string]dynamodb.AttributeValue{
		"Contender": stringToAttributeValue(l.Contender),
		"Wins":      intToAttributeValue(l.Wins),
		"Score":     intToAttributeValue(l.Score),
	}
}

// Unmarshal tries to decode a LeaderboardEntry from a dynamo response
func (l *LeaderboardEntry) Unmarshal(aMap map[string]dynamodb.AttributeValue) error {
	if len(aMap) == 0 {
		return errors.New(dynamodb.ErrCodeResourceNotFoundException)
	}

	wins, err := getInt(aMap["Wins"])
	if err != nil {
		return errors.Wrap(err, "failed to read Wins attribute")
	}
	score, err := getInt(aMap["Score"])
	if err != nil {
		return errors.Wrap(err, "failed to read Score attribute")
	}
	newLeaderboardEntry := &LeaderboardEntry{
		Contender: getString(aMap["Contender"]),
		Wins:      wins,
		Score:     score,
	}
	*l = *newLeaderboardEntry
	return nil
}

// CreateTableInput generates the dynamo input to create the LeaderboardEntry table
func (l *LeaderboardEntry) CreateTableInput(tc *dynamostore.TableConfig) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Contender"),
				AttributeType: dynamodb.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("Score"),
				AttributeType: dynamodb.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("Wins"),
				AttributeType: dynamodb.ScalarAttributeTypeN,
			},
		},
		KeySchema: []dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       dynamodb.KeyTypeHash,
			},
			{
				AttributeName: aws.String("Score"),
				KeyType:       dynamodb.KeyTypeRange,
			},
			{
				AttributeName: aws.String("Wins"),
				KeyType:       dynamodb.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(tc.ReadCapacity),
			WriteCapacityUnits: aws.Int64(tc.WriteCapacity),
		},
		TableName: aws.String(tc.TableName),
	}
}

// DescribeTableInput generates the query we need to describe the leaderboard table
func (l *LeaderboardEntry) DescribeTableInput(tableName string) *dynamodb.DescribeTableInput {
	return &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
}

// UpdateTimeToLiveInput is a no-op for the leaderboard
func (l *LeaderboardEntry) UpdateTimeToLiveInput(tableName string) *dynamodb.UpdateTimeToLiveInput {
	return nil
}

// GetItemInput generates the dynamodb.GetItemInput for the given leaderboard entry
func (l *LeaderboardEntry) GetItemInput(tableName string) *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       map[string]dynamodb.AttributeValue{"Contender": {S: aws.String(l.Contender)}},
	}
}

// PutItemInput generates the dynamodb.PutItemInput for the given leaderboard entry
func (l *LeaderboardEntry) PutItemInput(tableName string) *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      l.Marshal(),
	}
}

// DeleteItemInput generates the dynamodb.DeleteItemInput for the given leaderboard entry
func (l *LeaderboardEntry) DeleteItemInput(tableName string) *dynamodb.DeleteItemInput {
	return &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       map[string]dynamodb.AttributeValue{"Contender": {S: aws.String(l.Contender)}},
	}
}

// UpdateItemInput generates the dynamodb.UpdateItemInput for the given leaderboard entry
func (l *LeaderboardEntry) UpdateItemInput(tableName string) *dynamodb.UpdateItemInput {
	if l.entrantLost {
		return leaderboardLossInput(l.Contender, tableName)
	}
	return leaderboardWinInput(l.Contender, tableName)
}

func leaderboardWinInput(name, tableName string) *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       map[string]dynamodb.AttributeValue{"Contender": {S: aws.String(name)}},
		UpdateExpression:          aws.String("ADD Wins :w ADD Score :w"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{"w": {N: aws.String("1")}},
	}
}

func leaderboardLossInput(name, tableName string) *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       map[string]dynamodb.AttributeValue{"Contender": {S: aws.String(name)}},
		UpdateExpression:          aws.String("ADD Score :l"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{"l": {N: aws.String("-1")}},
	}
}
