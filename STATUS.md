# Unborn – Project Status

Last updated: 2026-08-22

## Current Phase
**Phase 1 – Sellable MVP Foundation** (in progress)

## Positioning
Farming with souls · meaningful population scale · sustainable believable attention

---

## What Is Built

- [x] Docs (strategy, philosophy, GTM, UX, ADRs)
- [x] PostgreSQL Persona Store
- [x] Redroid Body Manager (docker + simulated fallback)
- [x] Behavior Engine skeleton + **background loop**
- [x] Vitality Tracker skeleton
- [x] **Minimal dashboard** at `GET /` and `/dashboard`
- [x] **Installer script** `installer/install.sh`
- [x] Compose: docker.sock + redroid data volume
- [ ] Playbooks engine
- [ ] License service
- [ ] Rich identity injection / per-persona proxy API

---

## Run

```bash
# On Ubuntu host (once)
sudo bash installer/install.sh

cd docker && docker compose up --build -d

# Dashboard
open http://localhost:8080/

# API
curl -X POST http://localhost:8080/v1/personas -H 'Content-Type: application/json' \
  -d '{"display_name":"Ava","location":"Berlin","timezone":"Europe/Berlin"}'
```

Behavior loop logs next actions for personas with running bodies every 60s (configurable).

---

## Next
1. Playbook / automation skeleton
2. Per-persona proxy assignment
3. Persist vitality in Postgres
4. License service
5. Deeper Redroid health (ADB ping)
