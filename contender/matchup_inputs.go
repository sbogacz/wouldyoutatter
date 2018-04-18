package contender

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
)

const (
	matchupTableName = "Matchups"
)

// Key returns the Contenders name, and implements the dynamostore Item interface
func (m Matchup) Key() string {
	return m.Contender1 + m.Contender2
}

// Marshal encodes the values of a contender into the map format
// that dynamo expects
func (m Matchup) Marshal() map[string]dynamodb.AttributeValue {
	return map[string]dynamodb.AttributeValue{
		"Contender1":     stringToAttributeValue(m.Contender1),
		"Contender2":     stringToAttributeValue(m.Contender2),
		"Contender1Wins": intToAttributeValue(m.Contender1Wins),
		"Contender2Wins": intToAttributeValue(m.Contender2Wins),
	}
}

// Unmarshal tries to decode a Contender from a dynamo response
func (m *Matchup) Unmarshal(aMap map[string]dynamodb.AttributeValue) error {
	if len(aMap) == 0 {
		return errors.New(dynamodb.ErrCodeResourceNotFoundException)
	}

	contender1Wins, err := getInt(aMap["Contender1Wins"])
	if err != nil {
		return errors.Wrap(err, "failed to read Contender1Wins attribute")
	}
	contender2Wins, err := getInt(aMap["Contender2Wins"])
	if err != nil {
		return errors.Wrap(err, "failed to read Contender2Wins attribute")
	}
	newMatchup := &Matchup{
		Contender1:     getString(aMap["Contender1"]),
		Contender2:     getString(aMap["Contender2"]),
		Contender1Wins: contender1Wins,
		Contender2Wins: contender2Wins,
	}
	*m = *newMatchup
	return nil
}

// GetItemInput generates the dynamodb.GetItemInput for the given contender
func (m *Matchup) GetItemInput() *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(contenderTableName),
		Key: map[string]dynamodb.AttributeValue{
			"Contender1": {S: aws.String(m.Contender1)},
			"Contender2": {S: aws.String(m.Contender2)},
		},
	}
}

// PutItemInput generates the dynamodb.PutItemInput for the given contender
func (m *Matchup) PutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(contenderTableName),
		Item:      m.Marshal(),
	}
}

// DeleteItemInput generates the dynamodb.DeleteItemInput for the given contender
func (m *Matchup) DeleteItemInput() *dynamodb.DeleteItemInput {
	return &dynamodb.DeleteItemInput{
		TableName: aws.String(contenderTableName),
		Key: map[string]dynamodb.AttributeValue{
			"Contender1": {S: aws.String(m.Contender1)},
			"Contender2": {S: aws.String(m.Contender2)},
		},
	}
}

// UpdateItemInput generates the dynamodb.UpdateItemInput for the given contender
func (m *Matchup) UpdateItemInput() *dynamodb.UpdateItemInput {
	if m.contender1Won {
		return m.contender1WinInput()
	}
	return m.contender2WinInput()
}

func (m *Matchup) contender1WinInput() *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName: aws.String(matchupTableName),
		Key: map[string]dynamodb.AttributeValue{
			"Contender1": {S: aws.String(m.Contender1)},
			"Contender2": {S: aws.String(m.Contender2)},
		},
		UpdateExpression:          aws.String("ADD Contender1Wins :w ADD Contender2Losses :w"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{"w": {N: aws.String("1")}},
	}
}

func (m *Matchup) contender2WinInput() *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName: aws.String(matchupTableName),
		Key: map[string]dynamodb.AttributeValue{
			"Contender1": {S: aws.String(m.Contender1)},
			"Contender2": {S: aws.String(m.Contender2)},
		},
		UpdateExpression:          aws.String("ADD Contender2Wins :w ADD Contender1Losses :w"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{"w": {N: aws.String("1")}},
	}
}
