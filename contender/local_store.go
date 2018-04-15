package contender

import (
	"context"
	"sync"
)

type localStore struct {
	l          sync.RWMutex
	contenders map[string]*Contender
}

// NewLocalStore returns a local map backed instance of a Storer
func NewLocalStore() Storer {
	return &localStore{
		l:          sync.RWMutex{},
		contenders: make(map[string]*Contender, 10),
	}
}

func (s *localStore) Set(ctx context.Context, c Contender) error {
	s.l.Lock()
	s.contenders[c.Name] = &c
	s.l.Unlock()
	return nil
}

func (s *localStore) Get(ctx context.Context, name string) (*Contender, error) {
	s.l.RLock()
	c := s.contenders[name]
	s.l.RUnlock()
	return c, nil
}

func (s *localStore) Delete(ctx context.Context, name string) error {
	s.l.Lock()
	delete(s.contenders, name)
	s.l.Unlock()
	return nil
}
