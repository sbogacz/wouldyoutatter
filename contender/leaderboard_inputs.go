package contender

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
)

const (
	leaderboardTableName = "Leaderboard"
)

// Key returns the Leaderboards name, and implements the dynamostore Item interface
func (l Leaderboard) Key() string {
	return l.Contender
}

// Marshal encodes the values of a contender into the map format
// that dynamo expects
func (c Leaderboard) Marshal() map[string]dynamodb.AttributeValue {
	return map[string]dynamodb.AttributeValue{
		"Contender": stringToAttributeValue(c.Contender),
		"Wins":      intToAttributeValue(c.Wins),
		"Score":     intToAttributeValue(c.Score),
	}
}

// Unmarshal tries to decode a Leaderboard from a dynamo response
func (c *Leaderboard) Unmarshal(aMap map[string]dynamodb.AttributeValue) error {
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
	newLeaderboard := &Contender{
		Name:        getString(aMap["Name"]),
		Description: getString(aMap["Description"]),
		SVG:         getBytes(aMap["SVG"]),
		Wins:        wins,
		Losses:      losses,
		Score:       score,
	}
	*c = *newLeaderboard
	return nil
}

// GetItemInput generates the dynamodb.GetItemInput for the given contender
func (c *Leaderboard) GetItemInput() *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(contenderTableName),
		Key:       map[string]dynamodb.AttributeValue{"Name": {S: aws.String(c.Name)}},
	}
}

// PutItemInput generates the dynamodb.PutItemInput for the given contender
func (c *Leaderboard) PutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(contenderTableName),
		Item:      c.Marshal(),
	}
}

// DeleteItemInput generates the dynamodb.DeleteItemInput for the given contender
func (c *Leaderboard) DeleteItemInput() *dynamodb.DeleteItemInput {
	return &dynamodb.DeleteItemInput{
		TableName: aws.String(contenderTableName),
		Key:       map[string]dynamodb.AttributeValue{"Name": {S: aws.String(c.Name)}},
	}
}

// UpdateItemInput generates the dynamodb.UpdateItemInput for the given contender
func (c *Leaderboard) UpdateItemInput() *dynamodb.UpdateItemInput {
	if c.isLoser {
		return leaderboardLossInput(c.Name)
	}
	return leaderboardWinInput(c.Name)
}

func leaderboardWinInput(name string) *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName:                 aws.String(contenderTableName),
		Key:                       map[string]dynamodb.AttributeValue{"Name": {S: aws.String(name)}},
		UpdateExpression:          aws.String("ADD Wins :w"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{"w": {N: aws.String("1")}},
	}
}

func leaderboardLossInput(name string) *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName:                 aws.String(contenderTableName),
		Key:                       map[string]dynamodb.AttributeValue{"Name": {S: aws.String(name)}},
		UpdateExpression:          aws.String("ADD Losses :l"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{"l": {N: aws.String("1")}},
	}
}
