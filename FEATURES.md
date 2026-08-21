# Unborn – Features

Legend: [x] done · [ ] not yet

---

## Redroid / Execution

- [x] Privileged container start/stop
- [x] Per-body data volume + wipe API
- [x] ADB port allocation
- [x] Memory/CPU limits (env configurable)
- [x] Proxy boot props from persona assignment
- [x] Container + ADB health
- [x] Restart + logs APIs
- [x] Basic identity inject via ADB (model, manufacturer, android_id, …)
- [ ] Stronger persistent fingerprint
- [ ] Network namespace isolation
- [ ] Frida / touch injection

---

## Management plane

- [x] Persona store (Postgres + memory fallback)
- [x] Vitality (Postgres + memory fallback)
- [x] Behavior engine skeleton + background loop
- [x] Playbooks + assign
- [x] Per-persona proxy API
- [x] License service (offline HMAC)
- [x] Dashboard (population, bodies, playbooks, proxies, license)
- [x] Installer script
- [x] API reference doc
- [ ] Full playbook scheduler
- [ ] CLI

---

## Phase 2+

- [ ] Hierarchical memory, sensors, Radar, NL DNA
- [ ] Conversation + MCP
