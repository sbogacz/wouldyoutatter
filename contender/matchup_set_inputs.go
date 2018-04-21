package contender

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
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
	return map[string]dynamodb.AttributeValue{
		"ID": stringToAttributeValue(m.ID),
		"Set": {
			SS: m.Set,
		},
	}
}

// Unmarshal tries to decode a Contender from a dynamo response
func (m *MatchupSet) Unmarshal(aMap map[string]dynamodb.AttributeValue) error {
	set := []string{}
	setAttribute, ok := aMap["Set"]
	if !ok {
		return errors.New("no Set found")
	}
	if err := dynamodbattribute.Unmarshal(&setAttribute, &set); err != nil {
		return errors.Wrap(err, "failed to unmarshal Set")
	}
	newMatchupSet := &MatchupSet{
		ID:  getString(aMap["ID"]),
		Set: set,
	}
	*m = *newMatchupSet
	return nil
}

// CreateTableInput generates the dynamo input to create the matchupSet table
func (m *MatchupSet) CreateTableInput() *dynamodb.CreateTableInput {
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
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(m.tableName),
	}
}

// DescribeTableInput generates the query we need to describe the matchup set tables
func (m *MatchupSet) DescribeTableInput() *dynamodb.DescribeTableInput {
	return &dynamodb.DescribeTableInput{
		TableName: aws.String(m.tableName),
	}
}

// GetItemInput generates the dynamodb.GetItemInput for the given matchupSet
func (m *MatchupSet) GetItemInput() *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(m.tableName),
		Key: map[string]dynamodb.AttributeValue{
			"ID": {S: aws.String(m.ID)},
		},
	}
}

// PutItemInput generates the dynamodb.PutItemInput for the given matchupSet
func (m *MatchupSet) PutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(m.tableName),
		Item:      m.Marshal(),
	}
}

// DeleteItemInput generates the dynamodb.DeleteItemInput for the given matchupSet
func (m *MatchupSet) DeleteItemInput() *dynamodb.DeleteItemInput {
	return &dynamodb.DeleteItemInput{
		TableName: aws.String(m.tableName),
		Key: map[string]dynamodb.AttributeValue{
			"ID": {S: aws.String(m.ID)},
		},
	}
}

// UpdateItemInput is a no-op, since we don't update the matchupSet
func (m *MatchupSet) UpdateItemInput() *dynamodb.UpdateItemInput {
	updateExpression := "ADD Set :c"
	if m.entry.remove {
		updateExpression = "DELETE Set :c"
	}
	return &dynamodb.UpdateItemInput{
		TableName: aws.String(m.tableName),
		Key: map[string]dynamodb.AttributeValue{

			"ID": {S: aws.String(m.ID)},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{"c": {SS: []string{m.entry.String()}}},
	}
}

// Contenders is a collection that implements Scannable
type Contenders []Contender

// ScanInput producest a dynamodb ScanInput object
func (c *Contenders) ScanInput() *dynamodb.ScanInput {
	return &dynamodb.ScanInput{
		TableName: aws.String(contenderTableName),
	}
}

// Unmarshal allows results to be unmarshalled directly into the struct
func (c *Contenders) Unmarshal(maps []map[string]dynamodb.AttributeValue) error {
	cs := make([]*Contender, len(maps))
	for i := range cs {
		cs[i] = &Contender{}
		if err := cs[i].Unmarshal(maps[i]); err != nil {
			return errors.Wrap(err, "failed to unmarshal Contenders")
		}
	}
	contenders := make([]Contender, len(cs))
	for i := range cs {
		contenders[i] = *cs[i]
	}
	*c = contenders
	return nil

}
