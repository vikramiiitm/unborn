# Unborn – Project Status

Last updated: 2026-08-22

## Current Phase
**Phase 1 – Sellable MVP Foundation** (in progress)

## Positioning
Farming with souls. Meaningful population scale. Sustainable believable attention.

---

## What Is Built

### Documentation
- [x] Full strategy, philosophy, GTM, scale architecture, vitality design, ADRs

### Code
- [x] Persona Schema v0
- [x] Docker Compose (Postgres + Redis + Orchestrator)
- [x] PostgreSQL Persona Store (+ in-memory fallback)
- [x] Body Manager interface (simulated + Docker/Redroid skeleton)
- [x] Behavior Engine skeleton (rule-based, engagement-aware)
- [x] Vitality Tracker skeleton (0–100, levels, adjust API)
- [x] HTTP API:
  - personas, instances, device-profiles
  - `GET /v1/personas/{id}/next-action`
  - `GET /v1/personas/{id}/vitality` + `GET /v1/vitality`
- [ ] Real Redroid container start/stop (skeleton only)
- [ ] Behavior loop that actually drives bodies
- [ ] Vitality driven by Radar signals
- [ ] Installer / CLI / Dashboard

---

## Try it

```bash
cd docker && docker compose up --build

# Create persona
curl -X POST http://localhost:8080/v1/personas -H 'Content-Type: application/json' \
  -d '{"display_name":"Ava","location":"Berlin","timezone":"Europe/Berlin","engagement":"thoughtful_commenter"}'

# Next action for that persona
curl http://localhost:8080/v1/personas/<id>/next-action

# Vitality
curl http://localhost:8080/v1/personas/<id>/vitality
```

---

## Next
1. Wire real Redroid via Docker socket
2. Background behavior loop (persona → action → body)
3. Persist vitality + feed from future Radar
4. Installer script
