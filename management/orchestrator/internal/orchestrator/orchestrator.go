package orchestrator

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/config"
)

// InstanceState represents the lifecycle state of a body (Redroid instance).
type InstanceState string

const (
	InstanceStatePending   InstanceState = "pending"
	InstanceStateStarting  InstanceState = "starting"
	InstanceStateRunning   InstanceState = "running"
	InstanceStateStopping  InstanceState = "stopping"
	InstanceStateStopped   InstanceState = "stopped"
	InstanceStateFailed    InstanceState = "failed"
)

// Instance is a running (or planned) Redroid body bound to a Persona.
type Instance struct {
	ID           string        `json:"id"`
	PersonaID    string        `json:"persona_id"`
	State        InstanceState `json:"state"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	ContainerID  string        `json:"container_id,omitempty"`
	ErrorMessage string        `json:"error_message,omitempty"`
}

// Orchestrator is the core control plane for Personas and their bodies.
type Orchestrator struct {
	cfg       *config.Config
	mu        sync.RWMutex
	instances map[string]*Instance
}

func New(cfg *config.Config) *Orchestrator {
	return &Orchestrator{
		cfg:       cfg,
		instances: make(map[string]*Instance),
	}
}

// ListInstances returns all known instances.
func (o *Orchestrator) ListInstances() []*Instance {
	o.mu.RLock()
	defer o.mu.RUnlock()

	result := make([]*Instance, 0, len(o.instances))
	for _, inst := range o.instances {
		result = append(result, inst)
	}
	return result
}

// GetInstance returns a single instance by ID.
func (o *Orchestrator) GetInstance(id string) (*Instance, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	inst, ok := o.instances[id]
	return inst, ok
}

// CreateInstance registers a new instance bound to a Persona (skeleton – does not yet start Redroid).
func (o *Orchestrator) CreateInstance(personaID string) (*Instance, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if len(o.instances) >= o.cfg.MaxInstances {
		return nil, ErrMaxInstancesReached
	}

	now := time.Now().UTC()
	inst := &Instance{
		ID:        uuid.New().String(),
		PersonaID: personaID,
		State:     InstanceStatePending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	o.instances[inst.ID] = inst
	return inst, nil
}

// StopInstance marks an instance as stopping/stopped (skeleton).
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
