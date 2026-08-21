package orchestrator

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/config"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/identity"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/persona"
)

type InstanceState string

const (
	InstanceStatePending  InstanceState = "pending"
	InstanceStateStarting InstanceState = "starting"
	InstanceStateRunning  InstanceState = "running"
	InstanceStateStopping InstanceState = "stopping"
	InstanceStateStopped  InstanceState = "stopped"
	InstanceStateFailed   InstanceState = "failed"
)

type Instance struct {
	ID              string        `json:"id"`
	PersonaID       string        `json:"persona_id"`
	DeviceProfileID string        `json:"device_profile_id,omitempty"`
	State           InstanceState `json:"state"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	ContainerID     string        `json:"container_id,omitempty"`
	Simulated       bool          `json:"simulated"`
	ErrorMessage    string        `json:"error_message,omitempty"`
}

// PersonaRepository is the interface both memory and Postgres stores satisfy.
type PersonaRepository interface {
	Create(ctx context.Context, p *persona.Persona) (*persona.Persona, error)
	Get(ctx context.Context, id string) (*persona.Persona, error)
	List(ctx context.Context) ([]*persona.Persona, error)
	Update(ctx context.Context, p *persona.Persona) error
	Delete(ctx context.Context, id string) error
}

// memoryAdapter wraps the old in-memory store to the new interface.
type memoryAdapter struct {
	s *persona.Store
}

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
	cfg       *config.Config
	mu        sync.RWMutex
	instances map[string]*Instance
	personas  PersonaRepository
	profiles  map[string]*identity.DeviceProfile
}

func New(cfg *config.Config) *Orchestrator {
	ctx := context.Background()
	var repo PersonaRepository

	pg, err := persona.NewPostgresStore(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("postgres unavailable (%v) — falling back to in-memory Persona store", err)
		repo = &memoryAdapter{s: persona.NewStore()}
	} else {
		log.Println("Persona Store: PostgreSQL connected and migrated")
		repo = pg
	}

	o := &Orchestrator{
		cfg:       cfg,
		instances: make(map[string]*Instance),
		personas:  repo,
		profiles:  make(map[string]*identity.DeviceProfile),
	}

	for _, p := range identity.DefaultProfiles() {
		o.profiles[p.ID] = p
	}
	return o
}

func (o *Orchestrator) Personas() PersonaRepository {
	return o.personas
}

func (o *Orchestrator) ListInstances() []*Instance {
	o.mu.RLock()
	defer o.mu.RUnlock()
	result := make([]*Instance, 0, len(o.instances))
	for _, inst := range o.instances {
		result = append(result, inst)
	}
	return result
}

func (o *Orchestrator) GetInstance(id string) (*Instance, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	inst, ok := o.instances[id]
	return inst, ok
}

func (o *Orchestrator) CreateInstance(ctx context.Context, personaID string, simulated bool) (*Instance, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if len(o.instances) >= o.cfg.MaxInstances {
		return nil, ErrMaxInstancesReached
	}

	if _, err := o.personas.Get(ctx, personaID); err != nil {
		return nil, ErrPersonaNotFound
	}

	var profileID string
	for id := range o.profiles {
		profileID = id
		break
	}

	now := time.Now().UTC()
	inst := &Instance{
		ID:              uuid.New().String(),
		PersonaID:       personaID,
		DeviceProfileID: profileID,
		State:           InstanceStatePending,
		CreatedAt:       now,
		UpdatedAt:       now,
		Simulated:       simulated,
	}

	if simulated {
		inst.State = InstanceStateRunning
		inst.ContainerID = "sim-" + inst.ID[:8]
	}

	o.instances[inst.ID] = inst
	return inst, nil
}

func (o *Orchestrator) StopInstance(id string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	inst, ok := o.instances[id]
	if !ok {
		return ErrInstanceNotFound
	}
	inst.State = InstanceStateStopped
	inst.UpdatedAt = time.Now().UTC()
	return nil
}

func (o *Orchestrator) ListDeviceProfiles() []*identity.DeviceProfile {
	o.mu.RLock()
	defer o.mu.RUnlock()
	result := make([]*identity.DeviceProfile, 0, len(o.profiles))
	for _, p := range o.profiles {
		result = append(result, p)
	}
	return result
}
