package body

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type State string

const (
	StatePending  State = "pending"
	StateStarting State = "starting"
	StateRunning  State = "running"
	StateStopping State = "stopping"
	StateStopped  State = "stopped"
	StateFailed   State = "failed"
)

// NetworkOpts optional proxy for Redroid boot.
type NetworkOpts struct {
	ProxyHost string
	ProxyPort int
	ProxyType string // static, none
}

// StartOpts carries persona binding + network when creating a body.
type StartOpts struct {
	PersonaID       string
	DeviceProfileID string
	Simulated       bool
	Network         NetworkOpts
}

type Body struct {
	ID              string    `json:"id"`
	PersonaID       string    `json:"persona_id"`
	DeviceProfileID string    `json:"device_profile_id,omitempty"`
	State           State     `json:"state"`
	Simulated       bool      `json:"simulated"`
	ContainerID     string    `json:"container_id,omitempty"`
	ADBPort         int       `json:"adb_port,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ErrorMessage    string    `json:"error_message,omitempty"`
	Healthy         *bool     `json:"healthy,omitempty"`
}

type Manager interface {
	Start(ctx context.Context, opts StartOpts) (*Body, error)
	Stop(ctx context.Context, bodyID string) error
	Get(bodyID string) (*Body, bool)
	List() []*Body
}

type SimulatedManager struct {
	mu     sync.RWMutex
	bodies map[string]*Body
	max    int
}

func NewSimulatedManager(maxInstances int) *SimulatedManager {
	if maxInstances <= 0 {
		maxInstances = 10
	}
	return &SimulatedManager{bodies: make(map[string]*Body), max: maxInstances}
}

func (m *SimulatedManager) Start(ctx context.Context, opts StartOpts) (*Body, error) {
	_ = ctx
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.bodies) >= m.max {
		return nil, fmt.Errorf("max instances reached (%d)", m.max)
	}
	now := time.Now().UTC()
	h := true
	b := &Body{
		ID:              uuid.New().String(),
		PersonaID:       opts.PersonaID,
		DeviceProfileID: opts.DeviceProfileID,
		State:           StateRunning,
		Simulated:       true,
		ContainerID:     "sim-" + uuid.New().String()[:8],
		CreatedAt:       now,
		UpdatedAt:       now,
		Healthy:         &h,
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
