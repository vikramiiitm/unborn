# Unborn – Project Status

Last updated: 2026-08-22

## Current Phase
**Phase 1 – Sellable MVP Foundation** (in progress)

## Positioning (Locked)
We are a farming system. The moat is that every body has a Persona (soul).  
We operate at meaningful population scale. We optimize for sustainable, believable attention.

See `docs/farming-with-souls.md` and `docs/go-to-market.md`.

---

## What Is Built

### Documentation
- [x] Philosophy, positioning, GTM, scale architecture, vitality, ownership, stress test
- [x] Docs index, ADRs, tracking system

### Code
- [x] Repo structure
- [x] Persona Schema v0
- [x] Docker Compose foundation (Postgres + Redis + Orchestrator)
- [x] Orchestrator (Go)
  - **PostgreSQL-backed Persona Store** (with in-memory fallback)
  - Auto-migration on startup
  - DeviceProfile (basic identity)
  - Simulated body mode
  - HTTP API: health, personas, instances, device-profiles
- [ ] Real Redroid body lifecycle
- [ ] Behavior Engine
- [ ] Vitality module
- [ ] Installer, CLI, Dashboard

---

## How to try

```bash
cd docker
docker compose up --build

curl -X POST http://localhost:8080/v1/personas \
  -H "Content-Type: application/json" \
  -d '{"display_name":"Ava","location":"Berlin","timezone":"Europe/Berlin","age_min":26,"age_max":28}'

curl http://localhost:8080/v1/personas
```

Personas now persist in PostgreSQL across restarts.

---

## Next Technical Priorities
1. Real Redroid body lifecycle
2. Simple rule-based Behavior Engine
3. Vitality score skeleton
4. One-command installer
5. Network isolation + proxy forcing
