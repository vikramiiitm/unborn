# Unborn – Project Status

Last updated: 2026-08-22

## Redroid + UI

### Redroid
- docker run, limits, proxy, data, wipe, health
- **ADB identity inject** (auto on real start + `POST /v1/instances/{id}/inject-identity`)

### Dashboard (`/`)
Tabs wired to live APIs:
- **Population** — personas, vitality, Sim/Real body start
- **Bodies** — list, health, inject identity, stop
- **Playbooks** — list + assignments
- **Proxies** — list
- Header — license status + max instances

### Next Redroid
- Stronger fingerprint
- Network isolation
- Frida / input
