package orchestrator

import (
	"context"
	"log"
	"time"

	"github.com/vikramiiitm/unborn/management/orchestrator/internal/behavior"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/body"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/config"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/identity"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/license"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/persona"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/playbook"
	"github.com/vikramiiitm/unborn/management/orchestrator/internal/proxy"
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
	cfg       *config.Config
	personas  PersonaRepository
	bodies    body.Manager
	redroid   *body.RedroidManager
	behavior  *behavior.Engine
	vitality  vitality.Service
	playbooks *playbook.Store
	proxies   *proxy.Store
	license   *license.Service
	profiles  map[string]*identity.DeviceProfile
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

	var vit vitality.Service
	if pgt, err := vitality.NewPostgresTracker(ctx, cfg.DatabaseURL); err != nil {
		log.Printf("vitality postgres unavailable (%v) — memory", err)
		vit = vitality.NewMemoryService()
	} else {
		log.Println("Vitality: PostgreSQL connected")
		vit = vitality.NewPGService(pgt)
	}

	lic := license.NewService()
	max := cfg.MaxInstances
	if lim := lic.MaxInstances(); lim > 0 && lim < max {
		max = lim
	}

	rm := body.NewRedroidManager(max, body.DefaultRedroidConfig())
	o := &Orchestrator{
		cfg:       cfg,
		personas:  repo,
		bodies:    rm,
		redroid:   rm,
		behavior:  behavior.NewEngine(),
		vitality:  vit,
		playbooks: playbook.NewStore(),
		proxies:   proxy.NewStore(),
		license:   lic,
		profiles:  make(map[string]*identity.DeviceProfile),
	}
	for _, p := range identity.DefaultProfiles() {
		o.profiles[p.ID] = p
	}
	log.Printf("Redroid ready; license max=%d", max)
	return o
}

func (o *Orchestrator) Personas() PersonaRepository { return o.personas }
func (o *Orchestrator) Vitality() vitality.Service  { return o.vitality }
func (o *Orchestrator) Behavior() *behavior.Engine  { return o.behavior }
func (o *Orchestrator) Playbooks() *playbook.Store  { return o.playbooks }
func (o *Orchestrator) Proxies() *proxy.Store       { return o.proxies }
func (o *Orchestrator) License() *license.Service   { return o.license }

func (o *Orchestrator) ListInstances() []*body.Body { return o.bodies.List() }
func (o *Orchestrator) GetInstance(id string) (*body.Body, bool) {
	return o.bodies.Get(id)
}

func (o *Orchestrator) CreateInstance(ctx context.Context, personaID string, simulated bool) (*body.Body, error) {
	st := o.license.Status()
	if !st.Valid {
		return nil, ErrLicenseInvalid
	}
	if len(o.bodies.List()) >= st.MaxInstances {
		return nil, ErrMaxInstancesReached
	}
	if _, err := o.personas.Get(ctx, personaID); err != nil {
		return nil, ErrPersonaNotFound
	}
	var profileID string
	var profile *identity.DeviceProfile
	for id, p := range o.profiles {
		profileID = id
		profile = p
		break
	}
	opts := body.StartOpts{
		PersonaID:       personaID,
		DeviceProfileID: profileID,
		Simulated:       simulated,
		ExtraBootProps:  map[string]string{},
	}
	if px, ok := o.proxies.Get(personaID); ok {
		opts.Network = body.NetworkOpts{ProxyHost: px.Host, ProxyPort: px.Port, ProxyType: px.Type}
	}
	if profile != nil {
		opts.ExtraBootProps["unborn.device_model"] = profile.Model
		opts.ExtraBootProps["unborn.device_manufacturer"] = profile.Manufacturer
	}
	b, err := o.bodies.Start(ctx, opts)
	if err != nil {
		return nil, err
	}
	o.vitality.Ensure(personaID)
	if !b.Simulated && profile != nil && b.ADBPort > 0 {
		go func(port int, prof *identity.DeviceProfile, bodyID string) {
			ctx2, cancel := context.WithTimeout(context.Background(), 90*time.Second)
			defer cancel()
			if err := identity.InjectViaADB(ctx2, port, prof); err != nil {
				log.Printf("identity inject body=%s: %v", bodyID[:8], err)
			} else {
				log.Printf("identity inject ok body=%s model=%s", bodyID[:8], prof.Model)
			}
		}(b.ADBPort, profile, b.ID)
	}
	return b, nil
}

func (o *Orchestrator) StopInstance(ctx context.Context, id string) error {
	return o.bodies.Stop(ctx, id)
}

func (o *Orchestrator) RestartInstance(ctx context.Context, id string) error {
	if o.redroid == nil {
		return ErrInstanceNotFound
	}
	return o.redroid.Restart(ctx, id)
}

func (o *Orchestrator) InstanceLogs(ctx context.Context, id string, tail int) (string, error) {
	if o.redroid == nil {
		return "", ErrInstanceNotFound
	}
	return o.redroid.Logs(ctx, id, tail)
}

func (o *Orchestrator) WipeInstanceData(id string) error {
	if o.redroid == nil {
		return ErrInstanceNotFound
	}
	return o.redroid.WipeData(id)
}

func (o *Orchestrator) InjectIdentity(ctx context.Context, bodyID string) error {
	b, ok := o.bodies.Get(bodyID)
	if !ok {
		return ErrInstanceNotFound
	}
	if b.Simulated {
		return nil
	}
	prof := o.profiles[b.DeviceProfileID]
	if prof == nil {
		for _, p := range o.profiles {
			prof = p
			break
		}
	}
	if prof == nil || b.ADBPort <= 0 {
		return ErrInstanceNotFound
	}
	return identity.InjectViaADB(ctx, b.ADBPort, prof)
}

func (o *Orchestrator) Screenshot(ctx context.Context, bodyID string) ([]byte, error) {
	b, ok := o.bodies.Get(bodyID)
	if !ok {
		return nil, ErrInstanceNotFound
	}
	if b.Simulated {
		return nil, ErrSimulatedNoScreen
	}
	if b.ADBPort <= 0 {
		return nil, ErrInstanceNotFound
	}
	return body.ScreenshotPNG(ctx, b.ADBPort)
}

func (o *Orchestrator) ListDeviceProfiles() []*identity.DeviceProfile {
	out := make([]*identity.DeviceProfile, 0, len(o.profiles))
	for _, p := range o.profiles {
		out = append(out, p)
	}
	return out
}

func (o *Orchestrator) NextAction(ctx context.Context, personaID string) (behavior.Action, error) {
	p, err := o.personas.Get(ctx, personaID)
	if err != nil {
		return behavior.Action{}, ErrPersonaNotFound
	}
	act := o.behavior.NextAction(p, time.Now().UTC())
	for _, a := range o.playbooks.ListAssignments(personaID) {
		if !a.Active {
			continue
		}
		if pb, ok := o.playbooks.Get(a.PlaybookID); ok {
			act.Reason = act.Reason + " | playbook:" + pb.Name
			if pb.Kind == "warmup" && act.Type == "like" {
				act.Type = "scroll"
				act.Reason = "warmup: prefer observe over engage"
			}
			if pb.Kind == "presence" && pb.Params["window"] != "" {
				act.Metadata = map[string]string{"preferred_window": pb.Params["window"]}
			}
			break
		}
	}
	return act, nil
}

func (o *Orchestrator) CheckBodyHealth(ctx context.Context, bodyID string) (bool, string) {
	b, ok := o.bodies.Get(bodyID)
	if !ok {
		return false, "not found"
	}
	if b.Simulated {
		return true, "simulated"
	}
	ref := b.ContainerName
	if ref == "" {
		ref = b.ContainerID
	}
	if ref != "" && !body.ContainerRunning(ctx, ref) {
		return false, "container not running"
	}
	if b.ADBPort <= 0 {
		return false, "no adb port"
	}
	if body.CheckADB(ctx, b.ADBPort) {
		return true, "adb device"
	}
	if ref != "" && body.ContainerRunning(ctx, ref) {
		return false, "container up, adb not ready"
	}
	return false, "adb unreachable"
}
