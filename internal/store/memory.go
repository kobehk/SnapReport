package store

import (
	"sync"

	"SnapReport/internal/model"
)

type Store interface {
	Save(r model.Report)
	Get(id string) (model.Report, bool)
	List() []model.Report
}

type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]model.Report
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make(map[string]model.Report),
	}
}

func (s *MemoryStore) Save(r model.Report) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[r.ID] = r
}

func (s *MemoryStore) Get(id string) (model.Report, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.items[id]
	return r, ok
}

func (s *MemoryStore) List() []model.Report {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]model.Report, 0, len(s.items))
	for _, r := range s.items {
		out = append(out, r)
	}
	return out
}
