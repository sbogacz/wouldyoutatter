package dynamostore

import (
	"context"
	"sync"

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
	c      *TableConfig
	dynamo *dynamodb.DynamoDB
	lock   *sync.RWMutex
}

var _ Storer = (*dynamoStore)(nil)

// New takes a reference to a DynamoDB instance
// and returns the dynamo-backed version of the store
func New(d *dynamodb.DynamoDB, c *TableConfig) Storer {
	return &dynamoStore{
		c:      c,
		dynamo: d,
		lock:   &sync.RWMutex{},
	}
}

// Set takes a Item and tries to save it to Dynamo
func (s *dynamoStore) Set(ctx context.Context, item Item) error {
	var (
		numRetries int
		ok         bool
	)
	if numRetries, ok = ctx.Value(rKey).(int); !ok {
		numRetries = 0
	}
	ctx = context.WithValue(ctx, rKey, numRetries+1)

	if item.Key() == "" {
		return errors.New("must provide a non-empty name")
	}

	s.lock.RLock()
	input := item.PutItemInput(s.c.TableName)
	req := s.dynamo.PutItemRequest(input)

	if _, err := req.Send(); err != nil {
		s.lock.RUnlock()
		if createTableErr := s.createTableOnError(ctx, item, err); createTableErr != nil {
			log.WithField("num_retries", numRetries).WithError(err).Error("failed to set item")
			return errors.Wrapf(createTableErr, "failed to write Item %s to the database", item.Key())
		}
		// retry after creating table
		return s.Set(ctx, item)
	}
	s.lock.RUnlock()
	log.WithField("key", item.Key()).Debugf("successfully set item")

	return nil
}

// Get takes a name and tries to retrieve it from DynamoDB
func (s *dynamoStore) Get(ctx context.Context, item Item) (Item, error) {
	if item.Key() == "" {
		return nil, errors.New("must provide a non-empty name")
	}
	s.lock.RLock()
	defer s.lock.RUnlock()

	input := item.GetItemInput(s.c.TableName)
	req := s.dynamo.GetItemRequest(input)
	output, err := req.Send()
	if err != nil {
		return nil, errors.Wrap(err, "failed to send Get request")
	}

	return item, item.Unmarshal(output.Item)
}

// Update takes a Item and tries to update it in Dynamo
func (s *dynamoStore) Update(ctx context.Context, item Item) error {
	var (
		numRetries int
		ok         bool
	)
	if numRetries, ok = ctx.Value(rKey).(int); !ok {
		numRetries = 0
	}
	ctx = context.WithValue(ctx, rKey, numRetries+1)
	if item.Key() == "" {
		return errors.New("must provide a non-empty name")
	}

	s.lock.RLock()
	input := item.UpdateItemInput(s.c.TableName)
	req := s.dynamo.UpdateItemRequest(input)

	if _, err := req.Send(); err != nil {
		s.lock.RUnlock()
		if createTableErr := s.createTableOnError(ctx, item, err); createTableErr != nil {
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

	input := item.DeleteItemInput(s.c.TableName)
	req := s.dynamo.DeleteItemRequest(input)

	if _, err := req.Send(); err != nil {
		return errors.Wrap(err, "failed to send Delete request")
	}
	return nil
}

// Scan takes a scannable and tries to scan against DynamoDB
func (s *dynamoStore) Scan(ctx context.Context, items Scannable) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	input := items.ScanInput(s.c.TableName)
	req := s.dynamo.ScanRequest(input)
	output, err := req.Send()
	if err != nil {
		return errors.Wrap(err, "failed to send Scan request")
	}

	return items.Unmarshal(output.Items)
}

func (s *dynamoStore) createTableOnError(ctx context.Context, item Item, err error) error {
	if !TableNotFoundError(err) {
		log.WithError(err).Errorf("argh")
		return err
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	input := item.CreateTableInput(s.c)
	req := s.dynamo.CreateTableRequest(input)
	if _, err := req.Send(); err != nil {
		log.WithError(err).Errorf("failed to create table for %s", item.Key())
	}

	log.Debug("going to wait")
	describeInput := item.DescribeTableInput(s.c.TableName)
	if err := s.dynamo.WaitUntilTableExistsWithContext(ctx, describeInput); err != nil {
		log.WithError(err).Errorf("table for %s was not created in time", item.Key())
		return errors.Wrap(err, "table was not created in time")
	}
	log.Debug("done waiting")

	// enable TTL is so configured
	updateInput := item.UpdateTimeToLiveInput(s.c.TableName)
	if updateInput != nil {
		enableTTLReq := s.dynamo.UpdateTimeToLiveRequest(updateInput)
		if _, err := enableTTLReq.Send(); err != nil {
			log.WithError(err).Errorf("failed to enable TTL on %s, first try", s.c.TableName)
			if _, err := enableTTLReq.Send(); err != nil {
				log.WithError(err).Errorf("failed to enable TTL on %s, second try", s.c.TableName)
				return errors.Wrapf(err, "failed to enable TTL on %s, second try", s.c.TableName)
			}
		}
		log.Debug("going to wait to enable TTL")
		describeInput := item.DescribeTableInput(s.c.TableName)
		if err := s.dynamo.WaitUntilTableExistsWithContext(ctx, describeInput); err != nil {
			log.WithError(err).Errorf("table for %s didn't have TTL enabled in time", item.Key())
			return errors.Wrap(err, "TTL was not enabled in time")
		}
		log.Debug("done waiting for TTL to be enabeld")
	}
	return nil
}
