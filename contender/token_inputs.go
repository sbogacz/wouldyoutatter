package contender

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
)

const (
	tokenTableName = "Tokens"
)

var _ dynamostore.Item = (*Token)(nil)

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
	newToken := &Token{
		ID:         getString(aMap["ID"]),
		Contender1: getString(aMap["Contender1"]),
		Contender2: getString(aMap["Contender2"]),
	}
	*t = *newToken
	return nil
}

// CreateTableInput generates the dynamo input to create the token table
func (t *Token) CreateTableInput(tc *dynamostore.TableConfig) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: dynamodb.ScalarAttributeTypeS,
			},
		},
		KeySchema: []dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
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

// DescribeTableInput generates the query we need to describe the token table
func (t *Token) DescribeTableInput() *dynamodb.DescribeTableInput {
	return &dynamodb.DescribeTableInput{
		TableName: aws.String(tokenTableName),
	}
}

// UpdateTimeToLiveInput generates the input in order to set TTL on the token table
func (t *Token) UpdateTimeToLiveInput(tableName string) *dynamodb.UpdateTimeToLiveInput {
	return &dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String(tableName),
		TimeToLiveSpecification: &dynamodb.TimeToLiveSpecification{
			AttributeName: aws.String(s.c.TTLAttributeName),
			Enabled:       aws.Bool(true),
		},
	}
}

// GetItemInput generates the dynamodb.GetItemInput for the given token
func (t *Token) GetItemInput() *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(tokenTableName),
		Key: map[string]dynamodb.AttributeValue{
			"ID": {S: aws.String(t.ID)},
		},
	}
}

// PutItemInput generates the dynamodb.PutItemInput for the given token
func (t *Token) PutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(tokenTableName),
		Item:      t.Marshal(),
	}
}

// DeleteItemInput generates the dynamodb.DeleteItemInput for the given token
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
