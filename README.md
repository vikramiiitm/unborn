# Unborn

**Platform for deploying authentic digital people of specific niches and styles.**

Unborn turns customer-owned Ubuntu hardware into a self-hosted farm of high-fidelity, long-lived digital personas.  
These are not bots or generic device instances — they are coherent, believable people with demographics, taste, linguistic style, circadian rhythms, engagement philosophy, and memory.

Customers (brands, creators, agencies, communities) deploy the *right kind of people* so that real audiences perceive genuine niche members engaging with them.

---

## Core Idea

> You are not selling views, subscriptions, or device farms.  
> You are selling **specific types of people**.

The technical stack (Redroid, sensors, identity, network isolation) is infrastructure.  
**The moat is the quality, specificity, and believability of the personas.**

### Example Personas Customers Might Request
- 25–28 year old creative professionals in Berlin who care about design, slow living, and analog photography
- Early-30s fitness-oriented mothers in the US Midwest who engage with practical parenting + home workout content
- 18–22 year old streetwear and sneaker enthusiasts in Seoul / Tokyo style
- Quiet, thoughtful 40+ readers who comment meaningfully on long-form writing

---

## Product Framing

| Layer              | Typical tools              | Unborn                                      |
|--------------------|----------------------------|---------------------------------------------|
| Identity           | Random device fingerprint  | Coherent human identity (demographics + taste + history) |
| Behavior           | Generic scroll / like scripts | Style-specific, personality-driven engagement |
| Consistency        | Session-based              | Long-term persona memory and evolution      |
| Perception         | “Looks like a bot”         | “Feels like a real person from that world”  |
| Value to customer  | Cheap actions              | Credible social proof of the *right* audience |

---

## Architecture Overview

```
Ubuntu Host
├── CLI / Installer
├── Docker + required kernel modules
│
├── Management Plane
│   ├── Orchestrator + REST/gRPC API + Web Dashboard
│   ├── License Service
│   ├── Identity Manager
│   ├── Sensor & Input Provider
│   ├── Behavior Engine (hierarchical + AI)
│   ├── Learning & Adaptation Layer
│   └── Observability & Detection Radar
│
└── Execution Plane
    └── N × Isolated Redroid instances
        ├── Unique network path + proxy
        ├── Unique device identity
        ├── Sensor/touch injection hooks
        └── Optional on-device small models
```

Persona is the **first-class object**. Every Redroid instance is bound to a rich Persona.

---

## Key Capabilities (Roadmap)

### Phase 1 – Sellable MVP
- Multi-instance Redroid + resource control
- Per-instance proxy & basic identity
- Simple Orchestrator + license system
- One-command installer
- Basic statistical + rule-based behavior

### Phase 2 – Realism Leap
- High-quality sensor + touch injection
- Full Persona system
- Vision-based UI understanding
- Detection radar (basic)

### Phase 3 – Moat
- Self-improving behavior models
- Goal-oriented agents
- Digital twin personas with long-term memory
- Natural language Persona Studio
- Synthetic data flywheel
- Population-level realism (believable small communities)

---

## Repository Structure

```
unborn/
├── README.md
├── docs/
│   ├── architecture.md
│   ├── persona-model.md
│   ├── behavior-engine.md
│   └── roadmap.md
├── management/          # Orchestrator, API, dashboard, license, etc.
├── execution/           # Redroid related helpers, injection, sensors
├── personas/            # Persona schemas, generators, DNA tools
├── installer/           # One-command install scripts
├── docker/              # Compose files, Dockerfiles
└── experiments/         # Research & prototypes
```

---

## Status

Early stage. Architecture and product definition locked.  
Building Phase 1 scaffolding next.

---

*Unborn — digital people, born on your hardware.*
