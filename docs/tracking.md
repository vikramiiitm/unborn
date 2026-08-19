# How We Track What Is Built

## Three Layers

1. **GitHub Issues + Milestones + Labels**  
   Day-to-day work, bugs, epics, hard problems.

2. **Architecture Decision Records (`docs/adr/`)**  
   Permanent record of *why* significant technical choices were made.

3. **STATUS.md (root)**  
   High-level dashboard of what actually exists vs what is planned.  
   Updated in the same commit that lands meaningful work.

## Rules

- Every significant feature or subsystem that becomes real updates `STATUS.md`.
- Every architectural or stack decision gets an ADR.
- Issues are the source of truth for “what are we working on”.
- STATUS.md is the source of truth for “what already exists”.

This combination stays lightweight while giving clear visibility as the project grows.
