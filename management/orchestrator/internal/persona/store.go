package persona

import (
	"sync"
	"time"
)

// Store is the in-memory Persona store (kept for tests / fallback).
// Production uses PostgresStore.
type Store struct {
	mu       sync.RWMutex
	personas map[string]*Persona
}

func NewStore() *Store {
	return &Store{
		personas: make(map[string]*Persona),
	}
}

func (s *Store) Create(p *Persona) *Persona {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.ID == "" {
		p.ID = New("", "", "UTC", 25, 30, EngagementThoughtfulCommenter).ID
	}
	s.personas[p.ID] = p
	return p
}

func (s *Store) Get(id string) (*Persona, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.personas[id]
	return p, ok
}

func (s *Store) List() []*Persona {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Persona, 0, len(s.personas))
	for _, p := range s.personas {
		result = append(result, p)
	}
	return result
}

func (s *Store) Update(p *Persona) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.personas[p.ID]; !ok {
		return false
	}
	p.UpdatedAt = time.Now().UTC()
	s.personas[p.ID] = p
	return true
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.personas[id]; !ok {
		return false
	}
	delete(s.personas, id)
	return true
}
