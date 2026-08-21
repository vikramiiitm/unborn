package vitality

import (
	"sync"
	"time"
)

type Score struct {
	PersonaID  string    `json:"persona_id"`
	Value      float64   `json:"value"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastReason string    `json:"last_reason,omitempty"`
}

const (
	DefaultScore  = 75.0
	ThrivingMin   = 80.0
	StableMin     = 55.0
	UnderPressure = 30.0
	CriticalMin   = 10.0
)

func Level(v float64) string {
	switch {
	case v >= ThrivingMin:
		return "thriving"
	case v >= StableMin:
		return "stable"
	case v >= UnderPressure:
		return "under_pressure"
	case v >= CriticalMin:
		return "critical"
	default:
		return "collapsed"
	}
}

// Tracker is the in-memory implementation (fallback).
type Tracker struct {
	mu     sync.RWMutex
	scores map[string]*Score
}

func NewTracker() *Tracker {
	return &Tracker{scores: make(map[string]*Score)}
}

func (t *Tracker) Get(personaID string) *Score {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if s, ok := t.scores[personaID]; ok {
		return s
	}
	return &Score{PersonaID: personaID, Value: DefaultScore, UpdatedAt: time.Now().UTC(), LastReason: "default"}
}

func (t *Tracker) Ensure(personaID string) *Score {
	t.mu.Lock()
	defer t.mu.Unlock()
	if s, ok := t.scores[personaID]; ok {
		return s
	}
	s := &Score{PersonaID: personaID, Value: DefaultScore, UpdatedAt: time.Now().UTC(), LastReason: "initialized"}
	t.scores[personaID] = s
	return s
}

func (t *Tracker) Adjust(personaID string, delta float64, reason string) *Score {
	t.mu.Lock()
	defer t.mu.Unlock()
	s, ok := t.scores[personaID]
	if !ok {
		s = &Score{PersonaID: personaID, Value: DefaultScore}
		t.scores[personaID] = s
	}
	s.Value += delta
	if s.Value > 100 {
		s.Value = 100
	}
	if s.Value < 0 {
		s.Value = 0
	}
	s.UpdatedAt = time.Now().UTC()
	s.LastReason = reason
	return s
}

func (t *Tracker) List() []*Score {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]*Score, 0, len(t.scores))
	for _, s := range t.scores {
		out = append(out, s)
	}
	return out
}

// Service abstracts memory vs Postgres for the orchestrator.
type Service interface {
	Get(personaID string) *Score
	Ensure(personaID string) *Score
	Adjust(personaID string, delta float64, reason string) *Score
	List() []*Score
}

// MemoryService wraps Tracker.
type MemoryService struct{ t *Tracker }

func NewMemoryService() *MemoryService { return &MemoryService{t: NewTracker()} }

func (m *MemoryService) Get(id string) *Score                          { return m.t.Get(id) }
func (m *MemoryService) Ensure(id string) *Score                       { return m.t.Ensure(id) }
func (m *MemoryService) Adjust(id string, d float64, r string) *Score { return m.t.Adjust(id, d, r) }
func (m *MemoryService) List() []*Score                                { return m.t.List() }

// PGService wraps PostgresTracker with sync context helpers.
type PGService struct{ t *PostgresTracker }

func NewPGService(t *PostgresTracker) *PGService { return &PGService{t: t} }

func (p *PGService) Get(id string) *Score {
	s, err := p.t.Get(context.Background(), id)
	if err != nil {
		return &Score{PersonaID: id, Value: DefaultScore, LastReason: "error"}
	}
	return s
}
func (p *PGService) Ensure(id string) *Score {
	s, err := p.t.Ensure(context.Background(), id)
	if err != nil {
		return &Score{PersonaID: id, Value: DefaultScore, LastReason: "error"}
	}
	return s
}
func (p *PGService) Adjust(id string, d float64, r string) *Score {
	s, err := p.t.Adjust(context.Background(), id, d, r)
	if err != nil {
		return &Score{PersonaID: id, Value: DefaultScore, LastReason: "error"}
	}
	return s
}
func (p *PGService) List() []*Score {
	list, err := p.t.List(context.Background())
	if err != nil {
		return nil
	}
	return list
}
