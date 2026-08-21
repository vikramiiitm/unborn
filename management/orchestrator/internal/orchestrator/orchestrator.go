package orchestrator

import (
	"context"
	"log"
	"time"

	"github.com/vikramiiitm/unborn/management/orchestrator/internal/behavior"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/body"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/config"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/identity"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/persona"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/vitality"
)

type PersonaRepository interface {
	Create(ctx context.Context, p *persona.Persona) (*persona.Persona, error)
	Get(ctx context.Context, id string) (*persona.Persona, error)
	List(ctx context.Context) ([]*persona.Persona, error)
	Update(ctx context.Context, p *persona.Persona) error
	Delete(ctx context.Context, id string) error
}

type memoryAdapter struct{ s *persona.Store }

func (m *memoryAdapter) Create(ctx context.Context, p *persona.Persona) (*persona.Persona, error) {
	return m.s.Create(p), nil
}
func (m *memoryAdapter) Get(ctx context.Context, id string) (*persona.Persona, error) {
	p, ok := m.s.Get(id)
	if !ok {
		return nil, persona.ErrNotFound
	}
	return p, nil
}
func (m *memoryAdapter) List(ctx context.Context) ([]*persona.Persona, error) {
	return m.s.List(), nil
}
func (m *memoryAdapter) Update(ctx context.Context, p *persona.Persona) error {
	if !m.s.Update(p) {
		return persona.ErrNotFound
	}
	return nil
}
func (m *memoryAdapter) Delete(ctx context.Context, id string) error {
	if !m.s.Delete(id) {
		return persona.ErrNotFound
	}
	return nil
}

type Orchestrator struct {
	cfg      *config.Config
	personas PersonaRepository
	bodies   body.Manager
	behavior *behavior.Engine
	vitality *vitality.Tracker
	profiles map[string]*identity.DeviceProfile
}

func New(cfg *config.Config) *Orchestrator {
	ctx := context.Background()
	var repo PersonaRepository
	pg, err := persona.NewPostgresStore(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("postgres unavailable (%v) — in-memory Persona store", err)
		repo = &memoryAdapter{s: persona.NewStore()}
	} else {
		log.Println("Persona Store: PostgreSQL connected and migrated")
		repo = pg
	}

	o := &Orchestrator{
		cfg:      cfg,
		personas: repo,
		bodies:   body.NewDockerManager(cfg.MaxInstances, ""),
		behavior: behavior.NewEngine(),
		vitality: vitality.NewTracker(),
		profiles: make(map[string]*identity.DeviceProfile),
	}
	for _, p := range identity.DefaultProfiles() {
		o.profiles[p.ID] = p
	}
	return o
}

func (o *Orchestrator) Personas() PersonaRepository { return o.personas }
func (o *Orchestrator) Vitality() *vitality.Tracker  { return o.vitality }
func (o *Orchestrator) Behavior() *behavior.Engine   { return o.behavior }

func (o *Orchestrator) ListInstances() []*body.Body {
	return o.bodies.List()
}

func (o *Orchestrator) GetInstance(id string) (*body.Body, bool) {
	return o.bodies.Get(id)
}

func (o *Orchestrator) CreateInstance(ctx context.Context, personaID string, simulated bool) (*body.Body, error) {
	if _, err := o.personas.Get(ctx, personaID); err != nil {
		return nil, ErrPersonaNotFound
	}
	var profileID string
	for id := range o.profiles {
		profileID = id
		break
	}
	b, err := o.bodies.Start(ctx, personaID, profileID, simulated)
	if err != nil {
		return nil, err
	}
	o.vitality.Ensure(personaID)
	return b, nil
}

func (o *Orchestrator) StopInstance(ctx context.Context, id string) error {
	return o.bodies.Stop(ctx, id)
}

func (o *Orchestrator) ListDeviceProfiles() []*identity.DeviceProfile {
	out := make([]*identity.DeviceProfile, 0, len(o.profiles))
	for _, p := range o.profiles {
		out = append(out, p)
	}
	return out
}

// NextAction returns the next behavior action for a persona (Phase 1 rules).
func (o *Orchestrator) NextAction(ctx context.Context, personaID string) (behavior.Action, error) {
	p, err := o.personas.Get(ctx, personaID)
	if err != nil {
		return behavior.Action{}, ErrPersonaNotFound
	}
	return o.behavior.NextAction(p, time.Now().UTC()), nil
}
