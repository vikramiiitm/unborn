# ADR-0001: Record Architecture Decisions

## Status
Accepted

## Date
2026-08-20

## Context
As Unborn grows, we will make many technical decisions about stack, module boundaries, data models, AI layers, security, etc. Without a clear record, context is lost and the project drifts.

## Decision
We will use Architecture Decision Records (ADRs) stored in `docs/adr/` following a lightweight template.

Every significant decision that affects architecture, long-term optionality, module boundaries, or core tech choices must have an ADR.

## Consequences
Positive:
- Future contributors (including future us) understand why things are the way they are
- Easier to revisit or supersede decisions cleanly
- Forces clarity at decision time

Negative:
- Small process overhead

## Alternatives Considered
- Only using GitHub issues → too noisy and hard to find later
- Only putting decisions in chat/docs → gets lost
- Heavyweight architecture documentation → too slow for our pace
