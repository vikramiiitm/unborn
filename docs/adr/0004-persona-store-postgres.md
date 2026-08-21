# ADR-0004: Persona Store uses PostgreSQL

## Status
Accepted

## Date
2026-08-22

## Context
Personas must survive restarts, be queryable, and support customer export. In-memory store was only a scaffold.

## Decision
Primary Persona Store is PostgreSQL (JSONB document per persona + indexed metadata columns).  
In-memory store remains as automatic fallback when Postgres is unavailable (dev convenience).

Schema migrates automatically on Orchestrator startup.

## Consequences
Positive: persistence, easy backup/export path, aligns with ADR-0003.  
Negative: requires Postgres in the deployment (already in docker-compose).
