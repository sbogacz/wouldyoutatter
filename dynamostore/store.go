package dynamostore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Item interface {
	Key() string
	PutItemInput() *dynamodb.PutItemInput
	GetItemInput() *dynamodb.GetItemInput
	UpdateItemInput() *dynamodb.UpdateItemInput
	DeleteItemInput() *dynamodb.DeleteItemInput
	FromAttributeMap(map[string]dynamodb.AttributeValue) (Item, error)
}

// Storer is the interface to the K/V retrieval of Contenders
type Storer interface {
	Set(context.Context, Item) error
	Get(context.Context, Item) (Item, error)
	Update(context.Context, Item) error
	Delete(context.Context, Item) error
}
