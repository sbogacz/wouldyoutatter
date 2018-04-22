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
func (c *Contender) CreateTableInput(tc dynamostore.TableConfig) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Name"),
				AttributeType: dynamodb.ScalarAttributeTypeS,
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
		TableName: aws.String(tc.TableName),
	}
}

// DescribeTableInput generates the query we need to describe the contender table
func (c *Contender) DescribeTableInput() *dynamodb.DescribeTableInput {
	return &dynamodb.DescribeTableInput{
		TableName: aws.String(contenderTableName),
	}
}

// GetItemInput generates the dynamodb.GetItemInput for the given contender
func (c *Contender) GetItemInput() *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(contenderTableName),
		Key:       map[string]dynamodb.AttributeValue{"Name": {S: aws.String(c.Name)}},
	}
}

// PutItemInput generates the dynamodb.PutItemInput for the given contender
func (c *Contender) PutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(contenderTableName),
		Item:      c.Marshal(),
	}
}

// DeleteItemInput generates the dynamodb.DeleteItemInput for the given contender
func (c *Contender) DeleteItemInput() *dynamodb.DeleteItemInput {
	return &dynamodb.DeleteItemInput{
		TableName: aws.String(contenderTableName),
		Key:       map[string]dynamodb.AttributeValue{"Name": {S: aws.String(c.Name)}},
	}
}

// UpdateItemInput generates the dynamodb.UpdateItemInput for the given contender
func (c *Contender) UpdateItemInput() *dynamodb.UpdateItemInput {
	if c.isLoser {
		return lossInput(c.Name)
	}
	return winInput(c.Name)
}

func winInput(name string) *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName:                 aws.String(contenderTableName),
		Key:                       map[string]dynamodb.AttributeValue{"Name": {S: aws.String(name)}},
		UpdateExpression:          aws.String("ADD Wins :w"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{"w": {N: aws.String("1")}},
	}
}

func lossInput(name string) *dynamodb.UpdateItemInput {
	return &dynamodb.UpdateItemInput{
		TableName:                 aws.String(contenderTableName),
		Key:                       map[string]dynamodb.AttributeValue{"Name": {S: aws.String(name)}},
		UpdateExpression:          aws.String("ADD Losses :l"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{"l": {N: aws.String("1")}},
	}
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
