package dynamostore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
)

type dynamoStore struct {
	dynamo *dynamodb.DynamoDB
}

// New takes a reference to a DynamoDB instance
// and returns the dynamo-backed version of the store
func New(d *dynamodb.DynamoDB) Storer {
	return &dynamoStore{
		dynamo: d,
	}
}

// Set takes a Item and tries to save it to Dynamo
func (s *dynamoStore) Set(ctx context.Context, item Item) error {
	if item.Key() == "" {
		return errors.New("must provide a non-empty name")
	}
	input := item.PutItemInput()
	req := s.dynamo.PutItemRequest(input)

	if _, err := req.Send(); err != nil {
		if createTableErr := s.createTableOnError(ctx, item, err); err != nil {
			return errors.Wrapf(createTableErr, "failed to write Item %s to the database", item.Key())
		}
		// retry after creating table
		return s.Set(ctx, item)
	}

	return nil
}

// Get takes a name and tries to retrieve it from DynamoDB
func (s *dynamoStore) Get(ctx context.Context, item Item) (Item, error) {
	if item.Key() == "" {
		return nil, errors.New("must provide a non-empty name")
	}
	input := item.GetItemInput()
	req := s.dynamo.GetItemRequest(input)
	output, err := req.Send()
	if err != nil {
		return nil, errors.Wrap(err, "failed to send Get request")
	}

	return item, item.Unmarshal(output.Item)
}

// Update takes a Item and tries to update it in Dynamo
func (s *dynamoStore) Update(ctx context.Context, item Item) error {
	if item.Key() == "" {
		return errors.New("must provide a non-empty name")
	}
	input := item.UpdateItemInput()
	req := s.dynamo.UpdateItemRequest(input)

	if _, err := req.Send(); err != nil {
		if createTableErr := s.createTableOnError(ctx, item, err); err != nil {
			return errors.Wrap(createTableErr, "failed to send Update request")
		}
		// retry after creating table
		return s.Update(ctx, item)
	}
	return nil
}

// Delete takes a Item and tries to delete it from Dynamo
func (s *dynamoStore) Delete(ctx context.Context, item Item) error {
	if item.Key() == "" {
		return errors.New("must provide a non-empty name")
	}
	input := item.DeleteItemInput()
	req := s.dynamo.DeleteItemRequest(input)

	if _, err := req.Send(); err != nil {
		return errors.Wrap(err, "failed to send Delete request")
	}
	return nil
}

func (s *dynamoStore) createTableOnError(ctx context.Context, item Item, err error) error {
	if err.Error() != dynamodb.ErrCodeResourceNotFoundException {
		return err
	}

	input := item.CreateTableInput()
	req := s.dynamo.CreateTableRequest(input)
	if _, err := req.Send(); err != nil {
		return errors.Wrapf(err, "failed to create table for %s", item.Key())
	}

	return nil
}
