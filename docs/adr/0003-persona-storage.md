# ADR-0003: Persona Storage

## Status
Accepted

## Date
2026-08-20

## Context
Personas are the core object of the system. They must be versioned, serializable, exportable (customer ownership), and support hierarchical memory that can grow over long periods.

## Decision
- **Primary store**: PostgreSQL for structured Persona data + metadata
- **Memory blobs / large summaries**: Object storage or PostgreSQL JSONB (start with JSONB, move to object storage if needed)
- **Format**: Versioned JSON schema (Persona Schema v0 defined in `personas/schema/`)
- Every Persona is fully exportable as a single portable artifact

## Consequences
Positive:
- Strong querying and transactional guarantees
- Easy backup and customer export
- Clear versioning path
- Works fully self-hosted

Negative:
- Need good schema migration discipline
- JSONB memory will need careful summarization to stay bounded

## Alternatives Considered
- Pure file-based / SQLite — Too weak for concurrent access and querying at scale
- Fully embedded database only — Harder to operate and back up cleanly for customers
