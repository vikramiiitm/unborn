# Unborn – Project Status

Last updated: 2026-08-20

## Current Phase
**Phase 1 – Sellable MVP Foundation** (in progress)

## Tracking System

1. **GitHub Issues** — All work items, bugs, epics, hard problems
2. **Architecture Decision Records** (`docs/adr/`) — Why technical decisions were made
3. **STATUS.md** (this file) — High-level view of what is actually built
4. **FEATURES.md** — Feature-level checklist by phase

---

## What Is Built

### Documentation & Strategy
- [x] Full strategy, philosophy, ownership, scale architecture, tech stack, ADRs

### Code & Implementation
- [x] Repository structure
- [x] Persona Schema v0
- [x] FEATURES.md + STATUS.md tracking
- [x] Docker Compose foundation (Postgres + Redis + Orchestrator)
- [x] Orchestrator (Go)
  - Instance lifecycle with **simulated body mode**
  - Persona Store (in-memory for now)
  - DeviceProfile (basic identity)
  - HTTP API:
    - `GET /health`
    - `GET/POST /v1/personas`
    - `GET /v1/personas/{id}`
    - `GET/POST /v1/instances`
    - `GET /v1/instances/{id}`
    - `POST /v1/instances/{id}/stop`
    - `GET /v1/device-profiles`
- [ ] Real Redroid body management
- [ ] PostgreSQL-backed Persona Store
- [ ] Rich Identity generation
- [ ] Behavior Engine
- [ ] Installer
- [ ] CLI & Dashboard

---

## How to try the current skeleton

```bash
cd docker
docker compose up --build

# Create a Persona
curl -X POST http://localhost:8080/v1/personas \
  -H "Content-Type: application/json" \
  -d '{"display_name":"Ava","location":"Berlin","timezone":"Europe/Berlin","age_min":26,"age_max":28,"engagement":"thoughtful_commenter"}'

# Create a simulated instance for that persona
curl -X POST http://localhost:8080/v1/instances \
  -H "Content-Type: application/json" \
  -d '{"persona_id":"<id-from-previous-call>","simulated":true}'

# List everything
curl http://localhost:8080/v1/personas
curl http://localhost:8080/v1/instances
curl http://localhost:8080/v1/device-profiles
```

---

## Next Actions
1. PostgreSQL-backed Persona Store
2. Real Redroid container lifecycle (behind the same interface)
3. Simple rule-based Behavior Engine
4. Basic installer script
