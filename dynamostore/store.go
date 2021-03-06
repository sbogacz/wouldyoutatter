package dynamostore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Item is the interface for things we can store in using
// the Storer
type Item interface {
	Key() string
	PutItemInput(string) *dynamodb.PutItemInput
	GetItemInput(string) *dynamodb.GetItemInput
	UpdateItemInput(string) *dynamodb.UpdateItemInput
	DeleteItemInput(string) *dynamodb.DeleteItemInput
	CreateTableInput(c *TableConfig) *dynamodb.CreateTableInput
	DescribeTableInput(string) *dynamodb.DescribeTableInput
	TableOptions(string) []TableOption
	Marshal() map[string]dynamodb.AttributeValue
	Unmarshal(map[string]dynamodb.AttributeValue) error
}

// TableOption is an interface to specify requests that occur post-table
// creation, e.g. TTL enabling, or GSI creation
type TableOption interface {
	Send(db *dynamodb.DynamoDB) error
	Name() string
}

// Queryable is an interface for items whose tables can be queried
type Queryable interface {
	QueryInput(string, int) *dynamodb.QueryInput
	Unmarshal([]map[string]dynamodb.AttributeValue) error
}

// Scannable is an interface for items whose tables can be scanned
type Scannable interface {
	ScanInput(string) *dynamodb.ScanInput
	Unmarshal([]map[string]dynamodb.AttributeValue) error
}

// Storer is the interface to the K/V retrieval of Contenders
type Storer interface {
	Set(context.Context, Item) error
	Get(context.Context, Item) (Item, error)
	Update(context.Context, Item) error
	Delete(context.Context, Item) error
	Scan(context.Context, Scannable) error
	Query(context.Context, Queryable, int) error
}
