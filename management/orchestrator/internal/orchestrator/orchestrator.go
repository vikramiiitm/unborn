package orchestrator

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/config"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/identity"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/persona"
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

// Instance is a running (or planned) body bound to a Persona.
type Instance struct {
	ID              string        `json:"id"`
	PersonaID       string        `json:"persona_id"`
	DeviceProfileID string        `json:"device_profile_id,omitempty"`
	State           InstanceState `json:"state"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	ContainerID     string        `json:"container_id,omitempty"`
	Simulated       bool          `json:"simulated"` // true when running in simulated body mode
	ErrorMessage    string        `json:"error_message,omitempty"`
}

// Orchestrator is the core control plane for Personas and their bodies.
type Orchestrator struct {
	cfg       *config.Config
	mu        sync.RWMutex
	instances map[string]*Instance
	personas  *persona.Store
	profiles  map[string]*identity.DeviceProfile
}

func New(cfg *config.Config) *Orchestrator {
	o := &Orchestrator{
		cfg:       cfg,
		instances: make(map[string]*Instance),
		personas:  persona.NewStore(),
		profiles:  make(map[string]*identity.DeviceProfile),
	}

	// Seed some default device profiles
	for _, p := range identity.DefaultProfiles() {
		o.profiles[p.ID] = p
	}

	return o
}

// PersonaStore exposes the persona store.
func (o *Orchestrator) PersonaStore() *persona.Store {
	return o.personas
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

// CreateInstance registers a new instance bound to a Persona.
// In Phase 1 we support simulated bodies so the control plane can be developed
// without requiring a full Redroid environment.
func (o *Orchestrator) CreateInstance(personaID string, simulated bool) (*Instance, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if len(o.instances) >= o.cfg.MaxInstances {
		return nil, ErrMaxInstancesReached
	}

	// Ensure persona exists
	if _, ok := o.personas.Get(personaID); !ok {
		return nil, ErrPersonaNotFound
	}

	// Pick a simple device profile (later this will be smarter / persona-aware)
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
		// In simulated mode we immediately mark it running.
		inst.State = InstanceStateRunning
		inst.ContainerID = "sim-" + inst.ID[:8]
	}

	o.instances[inst.ID] = inst
	return inst, nil
}

// StopInstance marks an instance as stopped.
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

// ListDeviceProfiles returns available device profiles.
func (o *Orchestrator) ListDeviceProfiles() []*identity.DeviceProfile {
	o.mu.RLock()
	defer o.mu.RUnlock()

	result := make([]*identity.DeviceProfile, 0, len(o.profiles))
	for _, p := range o.profiles {
		result = append(result, p)
	}
	return result
}
