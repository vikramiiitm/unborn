# Management UX, Dashboard & Automations

What operators need to run Unborn as a **farm with souls** at meaningful scale.

Researched against existing device-farm UIs (DeviceFarm, OpenSTF-style hubs, FIRERPA, Farmly, etc.) and adapted for Persona + Vitality + population realism.

---

## 1. Jobs the UI Must Support

| Job | Why it matters |
|-----|----------------|
| See the whole living population at a glance | Scale = tens to hundreds of personas |
| Create / clone / archive Personas | Soul is the asset |
| Bind / unbind bodies (Redroid) | Bodies are disposable |
| Assign proxy + identity per body | Isolation |
| Watch Vitality and act on critical ones | Consequences system |
| Schedule presence / warm-up / campaigns | Farming output |
| Inspect one Persona deeply (memory, last actions) | Coherence debugging |
| Population health (not just single device) | Our differentiator |
| Automations without code | Agency operators |
| Audit log | Trust + multi-user later |

---

## 2. Core Screens (MVP → Later)

### A. Population Overview (home)
- Grid or table of Personas: name, niche tags, status (draft/warming/active/paused), Vitality score + level, bound body yes/no, last action, uptime
- Filters: status, vitality band, niche, simulated vs real body
- Bulk actions: pause, resume, assign proxy pool, run warm-up playbook
- Population health strip: % thriving / under pressure / critical, total active bodies, capacity used

### B. Persona Detail
- Full Persona DNA (demographics, engagement, circadian, interests)
- Memory summary (short + long-term later)
- Vitality history chart + last reasons
- Bound body: container id, ADB port, device profile, proxy
- Recent actions from Behavior Engine
- Actions: edit DNA, force next-action, pause soul, recycle body, export Persona

### C. Bodies / Farm
- List of Redroid (or simulated) instances
- Resource usage hints (when available)
- Start real vs simulated
- Stop / restart / wipe data volume
- Per-body proxy + network notes

### D. Automations / Playbooks
- Named recipes: “14-day warm-up”, “evening presence only”, “niche scroll + light engage”
- Attach playbook to one Persona or a population filter
- Schedule (cron-like or circadian-relative)
- Simple conditions: if vitality < 40 → pause body; if CAPTCHA spike → backoff (Radar later)

### E. Detection Radar (Phase 2)
- Warnings per persona / population
- Suggested adaptations
- Link into Vitality adjustments

### F. Settings
- Max instances, Redroid image, data root
- Default proxy pool
- License / capacity
- Export / backup Personas

---

## 3. Automations Operators Actually Need

From farm-software research + our model:

1. **Warm-up playbooks** — multi-day presence ramp, low risk first
2. **Circadian enforcement** — don’t force activity outside persona hours
3. **Proxy rotation policies** — sticky vs rotate on failure
4. **Bulk create from niche description** — NL DNA later; templates now
5. **Auto-pause on critical Vitality** — preserve soul, free body
6. **Health checks** — body up, ADB reachable, disk pressure
7. **Scheduled campaigns** — “this population focuses on X content window”
8. **Retry / backoff** — transient failures without killing coherence

Unborn-specific (competitors don’t have):
- Automations keyed on **Persona state**, not only device online/offline
- Population distinctiveness checks before scaling a niche
- Vitality as a first-class automation trigger

---

## 4. UX Principles for Unborn

1. **Soul-first navigation** — primary object is Persona, not container ID
2. **Population over single device** — health of the group is visible first
3. **Consequences visible** — Vitality never hidden
4. **Safe defaults** — simulated bodies easy; real Redroid explicit
5. **Agency-friendly** — bulk ops, playbooks, clear language (not only eng jargon)
6. **Auditability** — who changed what on which Persona

---

## 5. Phase Mapping

| Phase | UX scope |
|-------|----------|
| Phase 1 | Minimal API-driven ops + simple read-only dashboard (list personas, vitality, instances, next-action) |
| Phase 2 | Full Persona detail, playbooks, bulk actions, basic Radar view |
| Phase 3 | NL control, conversation with Persona, advanced population tools |

---

## 6. Minimal Dashboard Widgets (Phase 1 target)

- Capacity: active bodies / max
- Vitality distribution (counts per level)
- Recent personas created
- One-click: create persona → start simulated body
- Link to API docs / curl examples until UI is rich

---

## 7. What We Will Not Copy Blindly

Commodity farms optimize for:
- Screen grid of 50 phones
- Mass tap / mass follow scripts
- Price per action

We can offer multi-view later, but **our primary surface is population + Persona health**, not a wall of screens. Screen streaming is secondary (support/debug), not the product identity.
