package dynamostore

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type dynamoStore struct {
	dynamo *dynamodb.DynamoDB
	lock   *sync.RWMutex
}

// New takes a reference to a DynamoDB instance
// and returns the dynamo-backed version of the store
func New(d *dynamodb.DynamoDB) Storer {
	return &dynamoStore{
		dynamo: d,
		lock:   &sync.RWMutex{},
	}
}

// Set takes a Item and tries to save it to Dynamo
func (s *dynamoStore) Set(ctx context.Context, item Item) error {
	if item.Key() == "" {
		return errors.New("must provide a non-empty name")
	}

	s.lock.RLock()
	input := item.PutItemInput()
	req := s.dynamo.PutItemRequest(input)

	if _, err := req.Send(); err != nil {
		s.lock.RUnlock()
		if createTableErr := s.createTableOnError(ctx, item, err); err != nil {
			return errors.Wrapf(createTableErr, "failed to write Item %s to the database", item.Key())
		}
		// retry after creating table
		return s.Set(ctx, item)
	}
	s.lock.RUnlock()

	return nil
}

// Get takes a name and tries to retrieve it from DynamoDB
func (s *dynamoStore) Get(ctx context.Context, item Item) (Item, error) {
	if item.Key() == "" {
		return nil, errors.New("must provide a non-empty name")
	}
	s.lock.RLock()
	defer s.lock.RUnlock()

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

	s.lock.RLock()
	input := item.UpdateItemInput()
	req := s.dynamo.UpdateItemRequest(input)

	if _, err := req.Send(); err != nil {
		s.lock.RUnlock()
		if createTableErr := s.createTableOnError(ctx, item, err); err != nil {
			return errors.Wrap(createTableErr, "failed to send Update request")
		}
		// retry after creating table
		return s.Update(ctx, item)
	}
	s.lock.RUnlock()
	return nil
}

// Delete takes a Item and tries to delete it from Dynamo
func (s *dynamoStore) Delete(ctx context.Context, item Item) error {
	if item.Key() == "" {
		return errors.New("must provide a non-empty name")
	}
	s.lock.RLock()
	defer s.lock.RUnlock()

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
	s.lock.Lock()
	defer s.lock.Unlock()

	input := item.CreateTableInput()
	req := s.dynamo.CreateTableRequest(input)
	if _, err := req.Send(); err != nil {
		log.WithError(err).Errorf("failed to create table for %s", item.Key())
	}

	log.Info("going to wait")
	describeInput := item.DescribeTableInput()
	if err := s.dynamo.WaitUntilTableExistsWithContext(ctx, describeInput); err != nil {
		log.WithError(err).Errorf("table for %s was not created in time", item.Key())
		return errors.Wrap(err, "table was not created in time")
	}
	log.Info("done waiting")

	return nil
}
