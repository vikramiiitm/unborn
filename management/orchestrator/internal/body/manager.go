package body

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// State of a running body (Redroid instance or simulated).
type State string

const (
	StatePending  State = "pending"
	StateStarting State = "starting"
	StateRunning  State = "running"
	StateStopping State = "stopping"
	StateStopped  State = "stopped"
	StateFailed   State = "failed"
)

// Body is one execution unit bound to a Persona.
type Body struct {
	ID              string    `json:"id"`
	PersonaID       string    `json:"persona_id"`
	DeviceProfileID string    `json:"device_profile_id,omitempty"`
	State           State     `json:"state"`
	Simulated       bool      `json:"simulated"`
	ContainerID     string    `json:"container_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ErrorMessage    string    `json:"error_message,omitempty"`
}

// Manager starts and stops bodies. Simulated by default; real Redroid later.
type Manager interface {
	Start(ctx context.Context, personaID, deviceProfileID string, simulated bool) (*Body, error)
	Stop(ctx context.Context, bodyID string) error
	Get(bodyID string) (*Body, bool)
	List() []*Body
}

// SimulatedManager runs fully in-process (no Docker). Used for development.
type SimulatedManager struct {
	mu      sync.RWMutex
	bodies  map[string]*Body
	max     int
}

func NewSimulatedManager(maxInstances int) *SimulatedManager {
	if maxInstances <= 0 {
		maxInstances = 10
	}
	return &SimulatedManager{
		bodies: make(map[string]*Body),
		max:    maxInstances,
	}
}

func (m *SimulatedManager) Start(ctx context.Context, personaID, deviceProfileID string, simulated bool) (*Body, error) {
	_ = ctx
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.bodies) >= m.max {
		return nil, fmt.Errorf("max instances reached (%d)", m.max)
	}

	now := time.Now().UTC()
	b := &Body{
		ID:              uuid.New().String(),
		PersonaID:       personaID,
		DeviceProfileID: deviceProfileID,
		State:           StateRunning,
		Simulated:       true,
		ContainerID:     "sim-" + uuid.New().String()[:8],
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	m.bodies[b.ID] = b
	return b, nil
}

func (m *SimulatedManager) Stop(ctx context.Context, bodyID string) error {
	_ = ctx
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.bodies[bodyID]
	if !ok {
		return fmt.Errorf("body not found")
	}
	b.State = StateStopped
	b.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *SimulatedManager) Get(bodyID string) (*Body, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.bodies[bodyID]
	return b, ok
}

func (m *SimulatedManager) List() []*Body {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*Body, 0, len(m.bodies))
	for _, b := range m.bodies {
		out = append(out, b)
	}
	return out
}

// DockerManager is a skeleton for real Redroid containers.
// It currently falls back to simulated behavior until Docker socket + Redroid image are wired.
type DockerManager struct {
	inner      *SimulatedManager
	redroidImage string
}

func NewDockerManager(maxInstances int, redroidImage string) *DockerManager {
	if redroidImage == "" {
		redroidImage = "redroid/redroid:latest"
	}
	return &DockerManager{
		inner:        NewSimulatedManager(maxInstances),
		redroidImage: redroidImage,
	}
}

func (m *DockerManager) Start(ctx context.Context, personaID, deviceProfileID string, simulated bool) (*Body, error) {
	// Real Docker/Redroid path will be implemented here:
	// docker run --privileged -v ... redroid/redroid ...
	// For now always use simulated so the control plane keeps working.
	return m.inner.Start(ctx, personaID, deviceProfileID, true)
}

func (m *DockerManager) Stop(ctx context.Context, bodyID string) error {
	return m.inner.Stop(ctx, bodyID)
}

func (m *DockerManager) Get(bodyID string) (*Body, bool) {
	return m.inner.Get(bodyID)
}

func (m *DockerManager) List() []*Body {
	return m.inner.List()
}
