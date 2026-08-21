package playbook

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Playbook is a named automation recipe applied to personas or populations.
type Playbook struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Kind        string            `json:"kind"` // warmup, presence, campaign, custom
	Params      map[string]string `json:"params,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

// Assignment links a playbook to a persona.
type Assignment struct {
	ID         string    `json:"id"`
	PlaybookID string    `json:"playbook_id"`
	PersonaID  string    `json:"persona_id"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
}

// Store is an in-memory playbook registry (Phase 1).
type Store struct {
	mu          sync.RWMutex
	playbooks   map[string]*Playbook
	assignments map[string]*Assignment
}

func NewStore() *Store {
	s := &Store{
		playbooks:   make(map[string]*Playbook),
		assignments: make(map[string]*Assignment),
	}
	s.seed()
	return s
}

func (s *Store) seed() {
	now := time.Now().UTC()
	defaults := []*Playbook{
		{ID: uuid.New().String(), Name: "14-day warm-up", Description: "Low-risk presence ramp over two weeks", Kind: "warmup",
			Params: map[string]string{"days": "14", "intensity": "low"}, CreatedAt: now},
		{ID: uuid.New().String(), Name: "Evening presence", Description: "Activity biased to evening circadian window", Kind: "presence",
			Params: map[string]string{"window": "18:00-23:00"}, CreatedAt: now},
		{ID: uuid.New().String(), Name: "Niche scroll + light engage", Description: "Browse feed with selective engagement", Kind: "campaign",
			Params: map[string]string{"mode": "selective"}, CreatedAt: now},
	}
	for _, p := range defaults {
		s.playbooks[p.ID] = p
	}
}

func (s *Store) List() []*Playbook {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Playbook, 0, len(s.playbooks))
	for _, p := range s.playbooks {
		out = append(out, p)
	}
	return out
}

func (s *Store) Get(id string) (*Playbook, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.playbooks[id]
	return p, ok
}

func (s *Store) Create(name, description, kind string, params map[string]string) *Playbook {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := &Playbook{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Kind:        kind,
		Params:      params,
		CreatedAt:   time.Now().UTC(),
	}
	if p.Kind == "" {
		p.Kind = "custom"
	}
	s.playbooks[p.ID] = p
	return p
}

func (s *Store) Assign(playbookID, personaID string) (*Assignment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.playbooks[playbookID]; !ok {
		return nil, errPlaybookNotFound
	}
	a := &Assignment{
		ID:         uuid.New().String(),
		PlaybookID: playbookID,
		PersonaID:  personaID,
		Active:     true,
		CreatedAt:  time.Now().UTC(),
	}
	s.assignments[a.ID] = a
	return a, nil
}

func (s *Store) ListAssignments(personaID string) []*Assignment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Assignment, 0)
	for _, a := range s.assignments {
		if personaID == "" || a.PersonaID == personaID {
			out = append(out, a)
		}
	}
	return out
}

var errPlaybookNotFound = errString("playbook not found")

type errString string

func (e errString) Error() string { return string(e) }
