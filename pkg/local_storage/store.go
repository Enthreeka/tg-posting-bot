package store

import "sync"

type Store struct {
	store map[int64]*Data

	mu sync.RWMutex
}

type Data struct {
	Data          interface{}
	OperationType TypeCommand
	PreferMsgID   int
	CurrentMsgID  int
	ChannelID     int
}

func NewStore() *Store {
	return &Store{
		store: make(map[int64]*Data, 30),
	}
}

func (s *Store) Set(data *Data, userID int64) {
	s.Delete(userID)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[userID] = data
}

func (s *Store) Read(userID int64) (*Data, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	d, ok := s.store[userID]
	if !ok {
		return nil, false
	}

	return d, true
}

func (s *Store) Delete(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, userID)
}
