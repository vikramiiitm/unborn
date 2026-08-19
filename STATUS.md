# Unborn – Project Status

Last updated: 2026-08-20

## Current Phase
**Phase 1 – Sellable MVP Foundation** (in progress)

## Tracking System

1. **GitHub Issues** — All work items, bugs, epics, hard problems
2. **Architecture Decision Records** (`docs/adr/`) — Why technical choices were made
3. **STATUS.md** (this file) — High-level view of what is actually built
4. **FEATURES.md** — Feature-level checklist by phase

---

## What Is Built

### Documentation & Strategy
- [x] Product vision & positioning
- [x] Philosophy (full-power modular)
- [x] Ownership model (Hybrid)
- [x] Pre-build decisions
- [x] Long-term vision & future moats
- [x] Scale architecture
- [x] Tech stack direction
- [x] ADR process

### Code & Implementation
- [x] Repository structure
- [x] Persona Schema v0
- [x] FEATURES.md tracking
- [x] Docker Compose foundation (Postgres + Redis + Orchestrator)
- [x] Orchestrator skeleton (Go)
  - Basic instance lifecycle (create / list / get / stop)
  - Health endpoint
  - Max instances guard
  - Clean HTTP API skeleton
- [ ] Real Redroid body management
- [ ] Persona Store integration
- [ ] Identity Manager
- [ ] License service
- [ ] Behavior Engine skeleton
- [ ] Installer script
- [ ] Minimal CLI
- [ ] Minimal Dashboard

---

## Next Actions
1. Wire Orchestrator to real Docker/Redroid lifecycle
2. Implement Persona Store (Postgres)
3. Basic Identity profile generation
4. Simple rule-based Behavior Engine
5. One-command installer

---

## How to Update
Any significant piece of work that lands should update both `STATUS.md` and `FEATURES.md` in the same commit.
