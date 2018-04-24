package contender

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/sbogacz/wouldyoutatter/dynamostore"
)

var _ dynamostore.Item = (*MatchupSet)(nil)

// Key returns the Contenders name, and implements the dynamostore Item interface
func (m MatchupSet) Key() string {
	return m.ID
}

// Marshal encodes the values of a contender into the map format
// that dynamo expects
func (m MatchupSet) Marshal() map[string]dynamodb.AttributeValue {
	set := make([]string, 0, len(m.Set))
	for _, entry := range m.Set {
		set = append(set, entry.String())
	}
	return map[string]dynamodb.AttributeValue{
		"ID": stringToAttributeValue(m.ID),
		"MatchupSet": {
			SS: set,
		},
	}
}

// Unmarshal tries to decode a Contender from a dynamo response
func (m *MatchupSet) Unmarshal(aMap map[string]dynamodb.AttributeValue) error {
	setAttribute, ok := aMap["MatchupSet"]
	if !ok {
		return errors.New("no MatchupSet found")
	}
	set := make([]MatchupSetEntry, 0, len(setAttribute.SS))
	for _, entry := range setAttribute.SS {
		set = append(set, matchupEntryfromString(entry))
	}

	newMatchupSet := &MatchupSet{
		ID:  getString(aMap["ID"]),
		Set: set,
	}
	*m = *newMatchupSet
	return nil
}

// CreateTableInput generates the dynamo input to create the matchupSet table
func (m *MatchupSet) CreateTableInput(tc *dynamostore.TableConfig) *dynamodb.CreateTableInput {
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

// DescribeTableInput generates the query we need to describe the matchup set tables
func (m *MatchupSet) DescribeTableInput(tableName string) *dynamodb.DescribeTableInput {
	return &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
}

// TableOptions is a no-op for the matchup set table
func (m *MatchupSet) TableOptions(tableName string) []dynamostore.TableOption {
	return nil
}

// GetItemInput generates the dynamodb.GetItemInput for the given matchupSet
func (m *MatchupSet) GetItemInput(tableName string) *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]dynamodb.AttributeValue{
			"ID": {S: aws.String(m.ID)},
		},
		ConsistentRead: aws.Bool(true),
	}
}

// PutItemInput generates the dynamodb.PutItemInput for the given matchupSet
func (m *MatchupSet) PutItemInput(tableName string) *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      m.Marshal(),
	}
}

// DeleteItemInput generates the dynamodb.DeleteItemInput for the given matchupSet
func (m *MatchupSet) DeleteItemInput(tableName string) *dynamodb.DeleteItemInput {
	return &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]dynamodb.AttributeValue{
			"ID": {S: aws.String(m.ID)},
		},
	}
}

// UpdateItemInput is a no-op, since we don't update the matchupSet
func (m *MatchupSet) UpdateItemInput(tableName string) *dynamodb.UpdateItemInput {
	updateExpression := "ADD MatchupSet :c" //ADD MatchupSet :c"
	if m.entry.remove {
		updateExpression = "DELETE MatchupSet :c"
	}
	return &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]dynamodb.AttributeValue{

			"ID": {S: aws.String(m.ID)},
		},
		UpdateExpression: aws.String(updateExpression),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":c": {
				SS: []string{m.entry.String()},
			},
		},
	}
}
