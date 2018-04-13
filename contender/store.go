package contender

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
)

const (
	contenderTableName = "contenders"
)

// Store is the interface to the K/V retrieval of Contenders
type Store interface {
	Get(string) (*Contender, error)
	Set(Contender) error
	Delete(string) error
	Win(string) error
	Loss(string) error
}

type dynamoStore struct {
	dynamo *dynamodb.DynamoDB
}

// NewDynamoStore takes a reference to a DynamoDB instance
// and returns the dynamo-backed version of the store
func NewDynamoStore(d *dynamodb.DynamoDB) Store {
	return &dynamoStore{
		dynamo: d,
	}
}

// Get takes a name and tries to retrieve it from DynamoDB
func (s *dynamoStore) Get(name string) (*Contender, error) {
	if name == "" {
		return nil, errors.New("must provide a non-empty name")
	}
	input := getInput(name)
	req := s.dynamo.GetItemRequest(input)
	output, err := req.Send()
	if err != nil {
		return nil, errors.Wrap(err, "failed to send Get request")
	}
	return FromAttributeMap(output.Item)
}

// Set takes a Contender and tries to save it to Dynamo
func (s *dynamoStore) Set(c Contender) error {
	input := putInput(c)
	req := s.dynamo.PutItemRequest(input)

	if _, err := req.Send(); err != nil {
		return errors.Wrap(err, "failed to write Contender %s to the database")
	}

	return nil
}

// Delete takes a Contender and tries to delete it from Dynamo
func (s *dynamoStore) Delete(name string) error {
	if name == "" {
		return errors.New("must provide a non-empty name")
	}
	input := deleteInput(name)
	req := s.dynamo.DeleteItemRequest(input)

	if _, err := req.Send(); err != nil {
		return errors.Wrap(err, "failed to send Get request")
	}
	return nil
}

// Win records a win for the contender in the database
func (s *dynamoStore) Win(name string) error {
	if name == "" {
		return errors.New("must provide a non-empty name")
	}
	input := winInput(name)
	req := s.dynamo.UpdateItemRequest(input)

	if _, err := req.Send(); err != nil {
		return errors.Wrap(err, "failed to send Get request")
	}
	return nil
}

// Loss records a loss for the contender in the database
func (s *dynamoStore) Loss(name string) error {
	if name == "" {
		return errors.New("must provide a non-empty name")
	}
	input := lossInput(name)
	req := s.dynamo.UpdateItemRequest(input)

	if _, err := req.Send(); err != nil {
		return errors.Wrap(err, "failed to send Get request")
	}
	return nil
}

// Dynamo helpers
func getInput(name string) *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(contenderTableName),
		Key:       map[string]dynamodb.AttributeValue{"Name": dynamodb.AttributeValue{S: aws.String(name)}},
	}
}

func putInput(c Contender) *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(contenderTableName),
		Item:      c.ToAttributeMap(),
	}
}

func deleteInput(name string) *dynamodb.DeleteItemInput {
	return &dynamodb.DeleteItemInput{
		TableName: aws.String(contenderTableName),
		Key:       map[string]dynamodb.AttributeValue{"Name": dynamodb.AttributeValue{S: aws.String(name)}},
	}
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
