package dynamostore

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type retryKey string

const (
	rKey       retryKey = "retries"
	maxRetries          = 2
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
	var numRetries int
	var ok bool
	if numRetries, ok = ctx.Value(rKey).(int); !ok {
		numRetries = 0
	}
	ctx = context.WithValue(ctx, rKey, numRetries+1)

	if item.Key() == "" {
		return errors.New("must provide a non-empty name")
	}
	input := item.PutItemInput()

	req := s.dynamo.PutItemRequest(input)
	fmt.Printf("input: %+v\n\n", input)
	fmt.Printf("req: %+v\n\n", req)

	if _, err := req.Send(); err != nil {
		if createTableErr := s.createTableOnError(ctx, item, err); err != nil {
			log.WithError(err).Error("failed to set item")
			return errors.Wrapf(createTableErr, "failed to write Item %s to the database", item.Key())
		}
		// retry after creating table
		log.Infof("Retrying %d", numRetries)
		return s.Set(ctx, item)
	}

	//fmt.Printf("Output: %+v", output)
	return nil
}

// Get takes a name and tries to retrieve it from DynamoDB
func (s *dynamoStore) Get(ctx context.Context, item Item) (Item, error) {
	if item.Key() == "" {
		return nil, errors.New("must provide a non-empty name")
	}
	input := item.GetItemInput()
	fmt.Printf("input: %+v\n\n", input)
	req := s.dynamo.GetItemRequest(input)
	output, err := req.Send()
	if err != nil {
		return nil, errors.Wrap(err, "failed to send Get request")
	}
	fmt.Printf("output: %+v\n\n", output.Item)

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
	if !TableNotFoundError(err) {

		return err
	}
	fmt.Println("should be here")
	if ctx.Value(rKey).(int) >= maxRetries {
		return errors.New("hit max retries")
	}

	input := item.CreateTableInput()
	fmt.Printf("create input %+v\n\n", input)
	req := s.dynamo.CreateTableRequest(input)
	if _, createErr := req.Send(); createErr != nil {
		log.WithError(createErr).Error("failed to create table")
		//return errors.Wrapf(createErr, "failed to create table for %s", item.Key())
	}

	return nil
}
