# ADR-0002: Orchestrator Language

## Status
Accepted

## Date
2026-08-20

## Context
The Orchestrator is the heart of the Management Plane. It handles lifecycle of personas and Redroid instances, scheduling, health, resource accounting, and API. We need to choose the primary language.

## Decision
**Go** for the Orchestrator and core Management Plane services.

Python remains the language for Behavior Engine, Persona Intelligence, trajectory generation, and AI-heavy components.

## Consequences
Positive:
- Excellent concurrency and performance for managing many instances
- Single static binary → very easy distribution in the installer
- Strong isolation from the Python AI stack
- Good long-term maintainability for a control plane

Negative:
- Two primary languages in the project (Go + Python)
- Slightly higher initial development friction for some team members

## Alternatives Considered
- **Python (FastAPI)** — Faster to prototype, but weaker for long-running high-concurrency control plane and distribution
- **Rust** — Excellent performance and safety, but slower development velocity than we want for Phase 1
