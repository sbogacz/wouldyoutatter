package dynamostore

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

type localStore struct {
	l     sync.RWMutex
	items map[string]Item
}

// NewInMemoryStore returns a local map backed instance of a Storer
func NewInMemoryStore() Storer {
	return &localStore{
		l:     sync.RWMutex{},
		items: make(map[string]Item, 10),
	}
}

func (s *localStore) Set(ctx context.Context, item Item) error {
	if item.Key() == "" {
		return errors.New("must provide a non-empty name")
	}
	s.l.Lock()
	s.items[item.Key()] = item
	s.l.Unlock()
	return nil
}

func (s *localStore) Get(ctx context.Context, item Item) (Item, error) {
	if item.Key() == "" {
		return nil, errors.New("must provide a non-empty name")
	}
	s.l.RLock()
	it := s.items[item.Key()]
	s.l.RUnlock()
	return it, nil
}

func (s *localStore) Update(ctx context.Context, item Item) error {
	if item.Key() == "" {
		return errors.New("must provide a non-empty name")
	}
	s.l.Lock()
	s.items[item.Key()] = item
	s.l.Unlock()
	return nil
}

func (s *localStore) Delete(ctx context.Context, item Item) error {
	if item.Key() == "" {
		return errors.New("must provide a non-empty name")
	}
	s.l.Lock()
	delete(s.items, item.Key())
	s.l.Unlock()
	return nil
}
