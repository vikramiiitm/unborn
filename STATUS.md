# Unborn – Project Status

Last updated: 2026-08-22

## Phase 1 – Sellable MVP Foundation (strong progress)

### Built
- Persona + Vitality on Postgres
- Redroid manager with **per-persona proxy on docker run**
- Behavior loop with **playbook influence** + ADB health nudges
- Playbooks API + proxy API
- **License service** (HMAC offline keys, dev default, activate API)
- Minimal dashboard + installer

### API highlights
```
GET  /v1/license
POST /v1/license/activate  { "key": "..." }
GET  /v1/instances/{id}/health
PUT  /v1/personas/{id}/proxy  → applied on next real body start
```

### Next (Phase 1 polish)
- Richer playbook execution (schedules, intensity curves)
- Identity injection into Redroid
- CLI package
- Harden license secret management for production
