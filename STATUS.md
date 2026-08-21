# Unborn – Project Status

Last updated: 2026-08-22

## Current Phase
**Phase 1 – Sellable MVP Foundation** (in progress)

## What Is Built

- [x] Docs + positioning (farming with souls)
- [x] PostgreSQL Personas + **PostgreSQL Vitality**
- [x] Redroid body manager + simulated fallback
- [x] Behavior engine + background loop
- [x] Minimal dashboard (`/`)
- [x] Installer script
- [x] **Playbooks** (warmup / presence / campaign seeds + assign API)
- [x] **Per-persona proxy** API
- [ ] Wire proxy into Redroid `docker run` boot args per start
- [ ] License service
- [ ] Playbook execution (not only assign)
- [ ] ADB health checks

---

## New API

```
GET  /v1/playbooks
POST /v1/playbooks
POST /v1/playbooks/{id}/assign   { "persona_id": "..." }
GET  /v1/playbook-assignments?persona_id=

PUT  /v1/personas/{id}/proxy     { "host", "port", "type", "username", "password" }
GET  /v1/personas/{id}/proxy
DELETE /v1/personas/{id}/proxy
GET  /v1/proxies
```

Vitality now survives restarts when Postgres is up.

---

## Next
1. Apply persona proxy on real Redroid start
2. Execute playbooks in behavior loop
3. License service
4. ADB health
