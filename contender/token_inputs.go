package contender

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
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
		"ExpireAt":   int64ToAttributeValue(t.ExpireAt),
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
func (t *Token) DescribeTableInput(tableName string) *dynamodb.DescribeTableInput {
	return &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
}

// GetItemInput generates the dynamodb.GetItemInput for the given token
func (t *Token) GetItemInput(tableName string) *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]dynamodb.AttributeValue{
			"ID": {S: aws.String(t.ID)},
		},
	}
}

// PutItemInput generates the dynamodb.PutItemInput for the given token
func (t *Token) PutItemInput(tableName string) *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      t.Marshal(),
	}
}

// DeleteItemInput generates the dynamodb.DeleteItemInput for the given token
func (t *Token) DeleteItemInput(tableName string) *dynamodb.DeleteItemInput {
	return &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]dynamodb.AttributeValue{
			"ID": {S: aws.String(t.ID)},
		},
	}
}

// UpdateItemInput is a no-op, since we don't update the token
func (t *Token) UpdateItemInput(tableName string) *dynamodb.UpdateItemInput {
	return nil
}

// TableOptions returns the TTL table option the token store needs
func (t *Token) TableOptions(tableName string) []dynamostore.TableOption {
	input := t.updateTimeToLiveInput(tableName)
	return []dynamostore.TableOption{dynamostore.NewTTLOption(input)}
}

// updateTimeToLiveInput generates the input in order to set TTL on the token table
func (t *Token) updateTimeToLiveInput(tableName string) *dynamodb.UpdateTimeToLiveInput {
	return &dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String(tableName),
		TimeToLiveSpecification: &dynamodb.TimeToLiveSpecification{
			AttributeName: aws.String("ExpireAt"),
			Enabled:       aws.Bool(true),
		},
	}
}
