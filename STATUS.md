# Unborn – Project Status

Last updated: 2026-08-22

## Current Phase
**Phase 1 – Sellable MVP Foundation** (in progress)

## Positioning
Farming with souls. Meaningful population scale. Sustainable believable attention.

---

## What Is Built

### Documentation
- [x] Strategy, philosophy, GTM, scale architecture, vitality, farming-with-souls
- [x] **Management UX & Automations** (`docs/management-ux.md`)

### Code
- [x] PostgreSQL Persona Store
- [x] Behavior Engine skeleton
- [x] Vitality Tracker skeleton
- [x] **Redroid Body Manager** — real `docker run` path when Docker is available; simulated fallback
  - privileged container, data volume, ADB port allocation, optional proxy boot props
  - image default: `redroid/redroid:14.0.0_64only-latest`
- [x] HTTP API (personas, instances, next-action, vitality)
- [ ] Host kernel modules + installer checks (binder/ashmem)
- [ ] Continuous behavior loop
- [ ] Dashboard UI
- [ ] Playbooks / automations engine

---

## Redroid notes

Requires on host:
```bash
sudo apt install linux-modules-extra-$(uname -r)
sudo modprobe binder_linux devices="binder,hwbinder,vndbinder"
sudo modprobe ashmem_linux   # if needed
```

Orchestrator uses Docker CLI. Mount docker.sock in compose when ready for real bodies.

`POST /v1/instances` with `simulated: false` (and `?real=true`) attempts real Redroid when Docker works.

---

## Next
1. Compose: optional docker.sock + data volume for Redroid
2. Background behavior loop
3. Minimal dashboard (population + vitality)
4. Installer (kernel modules + docker + compose)
5. Playbook automation skeleton
