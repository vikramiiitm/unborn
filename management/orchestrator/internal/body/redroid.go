package body

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
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
	MemoryLimit      string
	CPULimit         string
}

func DefaultRedroidConfig() RedroidConfig {
	cfg := RedroidConfig{
		Image:       "redroid/redroid:14.0.0_64only-latest",
		DataRoot:    "/var/lib/unborn/redroid-data",
		BaseADBPort: 5555,
		Width:       1080,
		Height:      1920,
		DPI:         480,
		MemoryLimit: "3072m",
		CPULimit:    "2.0",
	}
	if v := os.Getenv("REDROID_IMAGE"); v != "" {
		cfg.Image = v
	}
	if v := os.Getenv("REDROID_DATA_ROOT"); v != "" {
		cfg.DataRoot = v
	}
	if v := os.Getenv("REDROID_MEMORY"); v != "" {
		cfg.MemoryLimit = v
	}
	if v := os.Getenv("REDROID_CPUS"); v != "" {
		cfg.CPULimit = v
	}
	if v := os.Getenv("REDROID_BASE_ADB_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.BaseADBPort = n
		}
	}
	return cfg
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

func isHostPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

func getUsedHostPorts() map[int]bool {
	used := make(map[int]bool)
	out, err := exec.Command("docker", "ps", "-a", "--format", "{{.Ports}}").CombinedOutput()
	if err != nil {
		return used
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if idx := strings.Index(part, "->"); idx != -1 {
				hostPart := part[:idx]
				if colonIdx := strings.LastIndex(hostPart, ":"); colonIdx != -1 {
					portStr := hostPart[colonIdx+1:]
					if p, err := strconv.Atoi(portStr); err == nil {
						used[p] = true
					}
				}
			}
		}
	}
	return used
}

func (m *RedroidManager) allocatePort() (int, error) {
	usedHostPorts := getUsedHostPorts()
	for p := m.cfg.BaseADBPort; p < m.cfg.BaseADBPort+m.max*2; p++ {
		if _, inMem := m.ports[p]; !inMem && !usedHostPorts[p] && isHostPortAvailable(p) {
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
	_ = os.MkdirAll(dataDir, 0o755)

	args := []string{
		"run", "-d", "--privileged",
		"--name", name,
		"-v", dataDir + ":/data",
		"-p", fmt.Sprintf("%d:5555", port),
		"--restart", "unless-stopped",
	}
	if m.cfg.MemoryLimit != "" {
		args = append(args, "--memory", m.cfg.MemoryLimit)
	}
	if m.cfg.CPULimit != "" {
		args = append(args, "--cpus", m.cfg.CPULimit)
	}
	args = append(args,
		m.cfg.Image,
		fmt.Sprintf("androidboot.redroid_width=%d", m.cfg.Width),
		fmt.Sprintf("androidboot.redroid_height=%d", m.cfg.Height),
		fmt.Sprintf("androidboot.redroid_dpi=%d", m.cfg.DPI),
	)

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
	for k, v := range opts.ExtraBootProps {
		if k != "" && v != "" {
			args = append(args, fmt.Sprintf("%s=%s", k, v))
		}
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		_ = exec.Command("docker", "rm", "-f", name).Run()
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
		ContainerName:   name,
		ADBPort:         port,
		DataDir:         dataDir,
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
	ref := b.ContainerName
	if ref == "" {
		ref = b.ContainerID
	}
	cmd := exec.CommandContext(ctx, "docker", "rm", "-f", ref)
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

func (m *RedroidManager) Restart(ctx context.Context, bodyID string) error {
	m.mu.Lock()
	b, ok := m.bodies[bodyID]
	m.mu.Unlock()
	if !ok || b.Simulated {
		return fmt.Errorf("restart only for real bodies")
	}
	ref := b.ContainerName
	if ref == "" {
		ref = b.ContainerID
	}
	out, err := exec.CommandContext(ctx, "docker", "restart", ref).CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker restart: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	m.mu.Lock()
	b.State = StateRunning
	b.UpdatedAt = time.Now().UTC()
	m.mu.Unlock()
	return nil
}

func (m *RedroidManager) Logs(ctx context.Context, bodyID string, tail int) (string, error) {
	m.mu.Lock()
	b, ok := m.bodies[bodyID]
	m.mu.Unlock()
	if !ok || b.Simulated {
		return "", fmt.Errorf("logs only for real bodies")
	}
	if tail <= 0 {
		tail = 100
	}
	ref := b.ContainerName
	if ref == "" {
		ref = b.ContainerID
	}
	out, err := exec.CommandContext(ctx, "docker", "logs", "--tail", strconv.Itoa(tail), ref).CombinedOutput()
	return string(out), err
}

func (m *RedroidManager) WipeData(bodyID string) error {
	m.mu.Lock()
	b, ok := m.bodies[bodyID]
	m.mu.Unlock()
	if !ok || b.DataDir == "" {
		return fmt.Errorf("no data dir")
	}
	if b.State == StateRunning {
		return fmt.Errorf("stop body before wipe")
	}
	return os.RemoveAll(b.DataDir)
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

func CheckADB(ctx context.Context, hostPort int) bool {
	if hostPort <= 0 {
		return false
	}
	addr := ADBAddr(hostPort)
	_ = exec.CommandContext(ctx, "adb", "connect", addr).Run()
	out, err := exec.CommandContext(ctx, "adb", "-s", addr, "get-state").CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "device")
}

func ContainerRunning(ctx context.Context, nameOrID string) bool {
	if nameOrID == "" {
		return false
	}
	out, err := exec.CommandContext(ctx, "docker", "inspect", "-f", "{{.State.Running}}", nameOrID).CombinedOutput()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}
