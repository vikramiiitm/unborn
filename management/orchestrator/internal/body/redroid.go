package body

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// RedroidConfig controls how real containers are launched.
type RedroidConfig struct {
	Image          string // e.g. redroid/redroid:14.0.0_64only-latest
	DataRoot       string // host path for per-instance /data mounts
	BaseADBPort    int    // first host port for ADB (5555+)
	Width          int
	Height         int
	DPI            int
	DefaultProxyHost string
	DefaultProxyPort int
}

func DefaultRedroidConfig() RedroidConfig {
	return RedroidConfig{
		Image:       "redroid/redroid:14.0.0_64only-latest",
		DataRoot:    "/var/lib/unborn/redroid-data",
		BaseADBPort: 5555,
		Width:       1080,
		Height:      1920,
		DPI:         480,
	}
}

// RedroidManager starts real Redroid containers via Docker CLI when available.
// Falls back to SimulatedManager if docker is missing or Start is called with simulated=true.
type RedroidManager struct {
	cfg    RedroidConfig
	mu     sync.Mutex
	bodies map[string]*Body
	ports  map[int]string // hostPort -> bodyID
	max    int
	sim    *SimulatedManager
}

func NewRedroidManager(maxInstances int, cfg RedroidConfig) *RedroidManager {
	if maxInstances <= 0 {
		maxInstances = 10
	}
	if cfg.Image == "" {
		cfg = DefaultRedroidConfig()
	}
	return &RedroidManager{
		cfg:    cfg,
		bodies: make(map[string]*Body),
		ports:  make(map[int]string),
		max:    maxInstances,
		sim:    NewSimulatedManager(maxInstances),
	}
}

func dockerAvailable() bool {
	cmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	return cmd.Run() == nil
}

func (m *RedroidManager) allocatePort() (int, error) {
	for p := m.cfg.BaseADBPort; p < m.cfg.BaseADBPort+m.max*2; p++ {
		if _, used := m.ports[p]; !used {
			return p, nil
		}
	}
	return 0, fmt.Errorf("no free ADB ports")
}

func (m *RedroidManager) Start(ctx context.Context, personaID, deviceProfileID string, simulated bool) (*Body, error) {
	if simulated || !dockerAvailable() {
		return m.sim.Start(ctx, personaID, deviceProfileID, true)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.bodies) >= m.max {
		return nil, fmt.Errorf("max instances reached (%d)", m.max)
	}

	port, err := m.allocatePort()
	if err != nil {
		return nil, err
	}

	bodyID := uuid.New().String()
	name := "unborn-" + bodyID[:8]
	dataDir := fmt.Sprintf("%s/%s", m.cfg.DataRoot, bodyID)

	args := []string{
		"run", "-d", "--privileged",
		"--name", name,
		"-v", dataDir + ":/data",
		"-p", fmt.Sprintf("%d:5555", port),
		"--restart", "unless-stopped",
		m.cfg.Image,
		fmt.Sprintf("androidboot.redroid_width=%d", m.cfg.Width),
		fmt.Sprintf("androidboot.redroid_height=%d", m.cfg.Height),
		fmt.Sprintf("androidboot.redroid_dpi=%d", m.cfg.DPI),
	}
	if m.cfg.DefaultProxyHost != "" {
		args = append(args,
			"androidboot.redroid_net_proxy_type=static",
			"androidboot.redroid_net_proxy_host="+m.cfg.DefaultProxyHost,
			fmt.Sprintf("androidboot.redroid_net_proxy_port=%d", m.cfg.DefaultProxyPort),
		)
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker run failed: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	containerID := strings.TrimSpace(string(out))
	if len(containerID) > 12 {
		containerID = containerID[:12]
	}

	now := time.Now().UTC()
	b := &Body{
		ID:              bodyID,
		PersonaID:       personaID,
		DeviceProfileID: deviceProfileID,
		State:           StateRunning,
		Simulated:       false,
		ContainerID:     containerID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	m.bodies[bodyID] = b
	m.ports[port] = bodyID
	return b, nil
}

func (m *RedroidManager) Stop(ctx context.Context, bodyID string) error {
	m.mu.Lock()
	b, ok := m.bodies[bodyID]
	m.mu.Unlock()
	if !ok {
		// may be simulated
		return m.sim.Stop(ctx, bodyID)
	}
	if b.Simulated || b.ContainerID == "" {
		return m.sim.Stop(ctx, bodyID)
	}

	cmd := exec.CommandContext(ctx, "docker", "rm", "-f", b.ContainerID)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("docker rm: %w (%s)", err, strings.TrimSpace(string(out)))
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	b.State = StateStopped
	b.UpdatedAt = time.Now().UTC()
	for p, id := range m.ports {
		if id == bodyID {
			delete(m.ports, p)
			break
		}
	}
	return nil
}

func (m *RedroidManager) Get(bodyID string) (*Body, bool) {
	m.mu.Lock()
	b, ok := m.bodies[bodyID]
	m.mu.Unlock()
	if ok {
		return b, true
	}
	return m.sim.Get(bodyID)
}

func (m *RedroidManager) List() []*Body {
	m.mu.Lock()
	out := make([]*Body, 0, len(m.bodies))
	for _, b := range m.bodies {
		out = append(out, b)
	}
	m.mu.Unlock()
	out = append(out, m.sim.List()...)
	return out
}

// ADBPort returns the host ADB port for a body if known.
func (m *RedroidManager) ADBPort(bodyID string) (int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for p, id := range m.ports {
		if id == bodyID {
			return p, true
		}
	}
	return 0, false
}

// ParsePort helper for tests/config.
func ParsePort(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
