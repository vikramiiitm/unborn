# Unborn

**Platform for deploying authentic digital people of specific niches and styles.**

Unborn turns customer-owned Ubuntu hardware into a self-hosted farm of high-fidelity, long-lived digital personas.  
These are not bots or generic device instances — they are coherent, believable people with demographics, taste, linguistic style, circadian rhythms, engagement philosophy, and memory.

---

## Current Status

See **[STATUS.md](STATUS.md)** for a live view of what is built and what is next.

We are in **Phase 1 – Sellable MVP Foundation**.

---

## Core Idea

> You are not selling views, subscriptions, or device farms.  
> You are selling **specific types of people**.

The technical stack (Redroid, sensors, identity, network isolation) is infrastructure.  
**The moat is the quality, specificity, and believability of the personas.**

---

## Key Documents

- [Philosophy](docs/philosophy.md) — Full-power modular engine, segregated responsibility
- [Product Strategy](docs/product-strategy.md) — 1-year moat and revenue focus
- [Ownership & Future Agents](docs/ownership-and-future-agents.md) — Hybrid ownership + conversational/MCP future
- [Scale Architecture](docs/scale-architecture.md) — How the system works at scale
- [Tech Stack](docs/tech-stack.md)
- [Pre-build Decisions](docs/pre-build-decisions.md)
- [Long-term Vision](docs/long-term-vision.md)
- [Architecture Decision Records](docs/adr/)

---

## Repository Structure

```
unborn/
├── STATUS.md                # What is actually built
├── docs/                    # Strategy, architecture, ADRs
├── personas/                # Schema, DNA, templates
├── management/              # Orchestrator, API, Persona Intelligence (Go + services)
├── execution/               # Redroid bodies, injection, isolation
├── docker/                  # Compose & Dockerfiles
├── installer/               # One-command installer
└── experiments/             # Research & prototypes
```

---

## Tracking

- **Issues** → what we are working on
- **ADRs** → why technical decisions were made
- **STATUS.md** → what already exists

---

*Unborn — digital people, born on your hardware.*
