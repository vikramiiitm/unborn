# Unborn – Features

Living list of features, their status, and phase.

Legend:
- [ ] Not started
- [~] In progress
- [x] Done (basic version)
- [+] Done and iterated

---

## Phase 1 – Sellable MVP

### Core Infrastructure
- [x] Repository structure & tracking system
- [x] Persona Schema v0
- [x] Architecture Decision Records process
- [x] FEATURES.md tracking
- [x] Docker Compose foundation
- [ ] One-command installer
- [ ] License service (signed keys, max instances)
- [ ] Resource limits & health checks (basic exists)

### Orchestrator
- [x] Orchestrator skeleton (Go)
- [x] Start / stop / list instances (simulated mode)
- [x] Bind Persona ↔ body (basic)
- [x] Max instances guard
- [ ] Real Redroid container lifecycle
- [ ] Health monitoring & auto-restart of bodies

### Personas
- [x] In-memory Persona Store
- [x] Create / list / get Persona API
- [ ] PostgreSQL-backed Persona Store
- [ ] Persona export (customer ownership)

### Identity & Network
- [x] Basic DeviceProfile structure
- [x] Default device profiles
- [ ] Rich identity generation
- [ ] Per-instance network isolation
- [ ] Forced proxy routing
- [ ] Identity injection into Redroid

### Behavior (Basic)
- [ ] Simple statistical / rule-based behavior engine
- [ ] Circadian-aware activity windows (basic)

### Interfaces
- [x] Minimal HTTP API (health, personas, instances, device-profiles)
- [ ] Minimal CLI
- [ ] Minimal Dashboard

---

## Phase 2 – Realism Leap

- [ ] Full Persona object lifecycle + PostgreSQL
- [ ] Hierarchical memory (short + long term)
- [ ] High-quality correlated sensor + touch trajectories
- [ ] Frida injection path
- [ ] Vision-based UI understanding (basic)
- [ ] Detection Radar (basic signals)
- [ ] Natural Language Persona DNA generator (first version)

---

## Phase 3 – Moat Features

- [ ] Long-term persona evolution
- [ ] Population-level realism controls
- [ ] Self-improving behavior signals (opt-in)
- [ ] Natural Language Persona Studio
- [ ] Conversational interface to Personas
- [ ] MCP tool-use capability
- [ ] Advanced Detection Radar + adaptation

---

## Cross-Cutting

- [x] Hybrid ownership model
- [x] Full-power modular philosophy
- [x] Scale architecture defined
- [ ] Human perception testing protocol
- [ ] Customer export of Personas
