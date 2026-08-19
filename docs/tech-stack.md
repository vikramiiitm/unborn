# Tech Stack (Initial Recommendations)

Aligned with: full-power modular design, self-hosted core, hybrid optional cloud, long-term optionality.

## Host & Runtime

| Layer | Choice | Reason |
|-------|--------|--------|
| Host OS | Ubuntu 24.04 LTS (or newer) | Stability + kernel module support for Redroid |
| Container | Docker + Docker Compose (Phase 1) → Kubernetes optional later | Simple for self-hosted customers |
| Android Runtime | Redroid | Best current open path for multi-instance Android on Linux |
| Isolation | Docker networks + network namespaces + proxy forcing | Strong isolation without extreme complexity |

## Management Plane

| Component | Recommended | Notes |
|-----------|-------------|-------|
| Orchestrator + API | Go or Python (FastAPI) | Go for performance & single binary, FastAPI for speed of development. Final choice in ADR. |
| Dashboard | React + Vite (or Next.js if needed) | Modern, fast, good enough |
| License | Signed JWT or Cryptlex / Keygen style | Offline-capable |
| Config / State | PostgreSQL or SQLite (small) + Redis for hot state | Start simple |
| Message / Events | NATS or Redis Streams | Lightweight |

## Persona & Behavior Intelligence

| Component | Recommended | Notes |
|-----------|-------------|-------|
| Persona object & memory store | PostgreSQL + object storage (or embedded) | Versioned, serializable Personas |
| Behavior Engine | Python | Best ecosystem for AI + rapid iteration |
| High-level Planner | Classical + small LLM / structured output | Keep deterministic core where possible |
| Trajectory & Sensors | Python + physics-informed models | Can later move hot paths to faster languages |
| Vision / UI understanding | Efficient VLM (local or optional cloud) | Phase 2+ |
| DNA Generation | Larger models (optional cloud or local heavy) | Hybrid as decided |

## AI Capability Layers

1. **Deterministic / Classical core** – scheduling, resource control, basic rules, identity injection
2. **Small specialized models** – trajectory generation, micro-timing, local decisions (can run near instances)
3. **Medium models** – UI understanding, short-term planning, comment style
4. **Large models (optional / central)** – Persona DNA from natural language, deep perception scoring, research experiments

We keep the most latency-sensitive and privacy-sensitive pieces runnable fully locally.

## Security of the Product Itself

### Principles
- Customer data (Personas, memory, logs) never leaves their hardware unless they explicitly opt in
- Least privilege between containers
- Signed license keys + instance limits enforced locally
- No hidden phone-home for core functionality
- Clear separation so a compromised dashboard cannot easily take over running personas
- Secrets (proxy credentials, license keys) handled properly (not in environment variables in plain sight)
- Audit logging of administrative actions

### Concrete measures (Phase 1 onward)
- Network policies between management and execution plane
- Read-only root filesystems where possible
- Regular scanning of our own images
- License service that works fully offline
- Customer-controlled encryption keys for sensitive Persona data (future)

## Language Summary (starting point)

- **Go or Python** for Orchestrator/API (decision via ADR)
- **Python** for Behavior / Persona / AI heavy parts
- **TypeScript/React** for Dashboard
- **Bash + Docker** for installer
- **Redroid** as the current body

Final choices for Orchestrator language and exact model stack will be recorded as Architecture Decision Records.
