package proxy

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Assignment is per-persona (or per-body) proxy configuration.
type Assignment struct {
	ID        string    `json:"id"`
	PersonaID string    `json:"persona_id"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Type      string    `json:"type"` // http, socks5, static (redroid boot)
	Username  string    `json:"username,omitempty"`
	// Password intentionally omitted from list responses later if needed
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store maps persona → proxy.
type Store struct {
	mu   sync.RWMutex
	byPersona map[string]*Assignment
}

func NewStore() *Store {
	return &Store{byPersona: make(map[string]*Assignment)}
}

func (s *Store) Set(personaID, host string, port int, typ, user, pass string) *Assignment {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	if typ == "" {
		typ = "http"
	}
	a, ok := s.byPersona[personaID]
	if !ok {
		a = &Assignment{ID: uuid.New().String(), PersonaID: personaID, CreatedAt: now}
		s.byPersona[personaID] = a
	}
	a.Host = host
	a.Port = port
	a.Type = typ
	a.Username = user
	a.Password = pass
	a.UpdatedAt = now
	return a
}

func (s *Store) Get(personaID string) (*Assignment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.byPersona[personaID]
	return a, ok
}

func (s *Store) Delete(personaID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byPersona[personaID]; !ok {
		return false
	}
	delete(s.byPersona, personaID)
	return true
}

func (s *Store) List() []*Assignment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Assignment, 0, len(s.byPersona))
	for _, a := range s.byPersona {
		out = append(out, a)
	}
	return out
}
