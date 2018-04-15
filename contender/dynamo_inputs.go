package contender

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
)

const (
	contenderTableName = "contenders"
)

// Key returns the Contenders name, and implements the dynamostore Item interface
func (c Contender) Key() string {
	return c.Name
}

// ToAttributeMap encodes the values of a contender into the map format
// that dynamo expects
func (c Contender) ToAttributeMap() map[string]dynamodb.AttributeValue {
	return map[string]dynamodb.AttributeValue{
		"Name":        stringToAttributeValue(c.Name),
		"Description": stringToAttributeValue(c.Description),
		"SVG":         bytesToAttributeValue(c.SVG),
		"Wins":        intToAttributeValue(c.Wins),
		"Losses":      intToAttributeValue(c.Wins),
		"Score":       intToAttributeValue(c.Wins),
	}
}

// FromAttributeMap tries to decode a Contender from a dynamo response
func (c *Contender) FromAttributeMap(aMap map[string]dynamodb.AttributeValue) (dynamostore.Item, error) {
	wins, err := getInt(aMap["Wins"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to read Wins attribute")
	}
	losses, err := getInt(aMap["Losses"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to read Losses attribute")
	}
	score, err := getInt(aMap["Score"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to read Score attribute")
	}

	return &Contender{
		Name:        getString(aMap["Name"]),
		Description: getString(aMap["Description"]),
		SVG:         getBytes(aMap["SVG"]),
		Wins:        wins,
		Losses:      losses,
		Score:       score,
	}, nil
}

// GetItemInput generates the dynamodb.GetItemInput for the given contender
func (c *Contender) GetItemInput() *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(contenderTableName),
		Key:       map[string]dynamodb.AttributeValue{"Name": dynamodb.AttributeValue{S: aws.String(c.Name)}},
	}
}

// PutItemInput generates the dynamodb.PutItemInput for the given contender
func (c *Contender) PutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(contenderTableName),
		Item:      c.ToAttributeMap(),
	}
}

// DeleteItemInput generates the dynamodb.DeleteItemInput for the given contender
func (c *Contender) DeleteItemInput() *dynamodb.DeleteItemInput {
	return &dynamodb.DeleteItemInput{
		TableName: aws.String(contenderTableName),
		Key:       map[string]dynamodb.AttributeValue{"Name": dynamodb.AttributeValue{S: aws.String(c.Name)}},
	}
}

// UpdateItemInput generates the dynamodb.UpdateItemInput for the given contender
func (c *Contender) UpdateItemInput() *dynamodb.UpdateItemInput {
	if c.isLoser {
		return lossInput(c.Name)
	}
	return winInput(c.Name)
}
func stringToAttributeValue(s string) dynamodb.AttributeValue {
	return dynamodb.AttributeValue{S: aws.String(s)}
}

func intToAttributeValue(n int) dynamodb.AttributeValue {
	return dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%d", n))}
}

func bytesToAttributeValue(b []byte) dynamodb.AttributeValue {
	return dynamodb.AttributeValue{B: b}
}

func getString(a dynamodb.AttributeValue) string {
	return *a.S
}

func getInt(a dynamodb.AttributeValue) (int, error) {
	return strconv.Atoi(*a.N)
}

func getBytes(a dynamodb.AttributeValue) []byte {
	return a.B
}

func winInput(name string) *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName:                 aws.String(contenderTableName),
		Key:                       map[string]dynamodb.AttributeValue{"Name": dynamodb.AttributeValue{S: aws.String(name)}},
		UpdateExpression:          aws.String("ADD Wins :w"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{"w": dynamodb.AttributeValue{N: aws.String("1")}},
	}
}

func lossInput(name string) *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName:                 aws.String(contenderTableName),
		Key:                       map[string]dynamodb.AttributeValue{"Name": dynamodb.AttributeValue{S: aws.String(name)}},
		UpdateExpression:          aws.String("ADD Losses :l"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{"l": dynamodb.AttributeValue{N: aws.String("1")}},
	}
}
