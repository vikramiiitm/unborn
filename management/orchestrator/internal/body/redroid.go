package body

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RedroidConfig struct {
	Image            string
	DataRoot         string
	BaseADBPort      int
	Width            int
	Height           int
	DPI              int
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

type RedroidManager struct {
	cfg    RedroidConfig
	mu     sync.Mutex
	bodies map[string]*Body
	ports  map[int]string
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
	return exec.Command("docker", "version", "--format", "{{.Server.Version}}").Run() == nil
}

func (m *RedroidManager) allocatePort() (int, error) {
	for p := m.cfg.BaseADBPort; p < m.cfg.BaseADBPort+m.max*2; p++ {
		if _, used := m.ports[p]; !used {
			return p, nil
		}
	}
	return 0, fmt.Errorf("no free ADB ports")
}

func (m *RedroidManager) Start(ctx context.Context, opts StartOpts) (*Body, error) {
	if opts.Simulated || !dockerAvailable() {
		return m.sim.Start(ctx, opts)
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

	// Per-start proxy: persona assignment wins over default config
	proxyHost := opts.Network.ProxyHost
	proxyPort := opts.Network.ProxyPort
	if proxyHost == "" {
		proxyHost = m.cfg.DefaultProxyHost
		proxyPort = m.cfg.DefaultProxyPort
	}
	if proxyHost != "" {
		args = append(args,
			"androidboot.redroid_net_proxy_type=static",
			"androidboot.redroid_net_proxy_host="+proxyHost,
			fmt.Sprintf("androidboot.redroid_net_proxy_port=%d", proxyPort),
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
		PersonaID:       opts.PersonaID,
		DeviceProfileID: opts.DeviceProfileID,
		State:           StateRunning,
		Simulated:       false,
		ContainerID:     containerID,
		ADBPort:         port,
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
	return append(out, m.sim.List()...)
}

func (m *RedroidManager) ADBPort(bodyID string) (int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if b, ok := m.bodies[bodyID]; ok && b.ADBPort > 0 {
		return b.ADBPort, true
	}
	for p, id := range m.ports {
		if id == bodyID {
			return p, true
		}
	}
	return 0, false
}

// CheckADB pings adb connect to host:port (best-effort).
func CheckADB(ctx context.Context, hostPort int) bool {
	if hostPort <= 0 {
		return false
	}
	addr := fmt.Sprintf("127.0.0.1:%d", hostPort)
	cmd := exec.CommandContext(ctx, "adb", "connect", addr)
	if err := cmd.Run(); err != nil {
		// adb may be missing — try docker exec fallback is out of scope; mark unknown
		return false
	}
	cmd = exec.CommandContext(ctx, "adb", "-s", addr, "get-state")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "device")
}
