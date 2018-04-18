package contender

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
)

const (
	tokenTableName = "Tokens"
)

// Key returns the Contenders name, and implements the dynamostore Item interface
func (t Token) Key() string {
	return t.ID
}

// Marshal encodes the values of a contender into the map format
// that dynamo expects
func (t Token) Marshal() map[string]dynamodb.AttributeValue {
	return map[string]dynamodb.AttributeValue{
		"ID":         stringToAttributeValue(t.ID),
		"Contender1": stringToAttributeValue(t.Contender1),
		"Contender2": stringToAttributeValue(t.Contender2),
	}
}

// Unmarshal tries to decode a Contender from a dynamo response
func (t *Token) Unmarshal(aMap map[string]dynamodb.AttributeValue) error {
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
	newToken := &Token{
		ID:         getString(aMap["ID"]),
		Contender1: getString(aMap["Contender1"]),
		Contender2: getString(aMap["Contender2"]),
	}
	*t = *newToken
	return nil
}

// GetItemInput generates the dynamodb.GetItemInput for the given contender
func (t *Token) GetItemInput() *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(tokenTableName),
		Key: map[string]dynamodb.AttributeValue{
			"ID": {S: aws.String(t.ID)},
		},
	}
}

// PutItemInput generates the dynamodb.PutItemInput for the given contender
func (t *Token) PutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(tokenTableName),
		Item:      t.Marshal(),
	}
}

// DeleteItemInput generates the dynamodb.DeleteItemInput for the given contender
func (t *Token) DeleteItemInput() *dynamodb.DeleteItemInput {
	return &dynamodb.DeleteItemInput{
		TableName: aws.String(tokenTableName),
		Key: map[string]dynamodb.AttributeValue{
			"ID": {S: aws.String(t.ID)},
		},
	}
}

// UpdateItemInput is a no-op, since we don't update the token
func (t *Token) UpdateItemInput() *dynamodb.UpdateItemInput {
	return nil
}
