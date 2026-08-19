# Unborn – Project Status

Last updated: 2026-08-20

## Current Phase
**Phase 1 – Sellable MVP Foundation** (in progress)

## Positioning (Locked)
We are a farming system. The moat is that every body has a Persona (soul).  
We operate at meaningful population scale. We optimize for sustainable, believable attention.

See `docs/farming-with-souls.md` and `docs/go-to-market.md`.

---

## What Is Built

### Documentation (Cleaned)
- [x] Philosophy, positioning, GTM, scale architecture, vitality, ownership, stress test
- [x] Docs index (`docs/README.md`)
- [x] ADRs + tracking system

### Code
- [x] Repo structure
- [x] Persona Schema v0
- [x] Docker Compose foundation (Postgres + Redis + Orchestrator)
- [x] Orchestrator (Go)
  - Persona Store (in-memory)
  - DeviceProfile (basic identity)
  - Simulated body mode
  - HTTP API for personas, instances, device profiles
- [ ] PostgreSQL-backed Persona Store
- [ ] Real Redroid lifecycle
- [ ] Behavior Engine
- [ ] Vitality module
- [ ] Installer, CLI, Dashboard

---

## How to try the current skeleton

```bash
cd docker
docker compose up --build

curl -X POST http://localhost:8080/v1/personas \
  -H "Content-Type: application/json" \
  -d '{"display_name":"Ava","location":"Berlin","timezone":"Europe/Berlin","age_min":26,"age_max":28}'

curl -X POST http://localhost:8080/v1/instances \
  -H "Content-Type: application/json" \
  -d '{"persona_id":"<id>","simulated":true}'
```

---

## Next Technical Priorities
1. PostgreSQL Persona Store
2. Real Redroid body lifecycle
3. Simple Behavior Engine
4. Vitality score skeleton
5. Installer
