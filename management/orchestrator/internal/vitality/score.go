package vitality

import (
	"sync"
	"time"
)

// Score is 0–100. Higher = healthier, more coherent, surviving better.
type Score struct {
	PersonaID   string    `json:"persona_id"`
	Value       float64   `json:"value"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastReason  string    `json:"last_reason,omitempty"`
}

const (
	DefaultScore   = 75.0
	ThrivingMin    = 80.0
	StableMin      = 55.0
	UnderPressure  = 30.0
	CriticalMin    = 10.0
)

// Level derives consequence band from score.
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

// Tracker holds vitality scores in memory (Phase 1). Later: persist + Radar inputs.
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
	return &Score{
		PersonaID: personaID,
		Value:     DefaultScore,
		UpdatedAt: time.Now().UTC(),
		LastReason: "default",
	}
}

func (t *Tracker) Ensure(personaID string) *Score {
	t.mu.Lock()
	defer t.mu.Unlock()
	if s, ok := t.scores[personaID]; ok {
		return s
	}
	s := &Score{
		PersonaID:  personaID,
		Value:      DefaultScore,
		UpdatedAt:  time.Now().UTC(),
		LastReason: "initialized",
	}
	t.scores[personaID] = s
	return s
}

// Adjust applies a delta with a reason. Clamps to 0–100.
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
