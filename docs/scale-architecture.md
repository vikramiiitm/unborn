# Unborn Scale Architecture

This document describes how the system works at different scales, what happens in key scenarios, and what is required for each part — while respecting every decision we have locked.

---

## 1. Core Architectural Planes (Recap + Scale View)

```
┌─────────────────────────────────────────────────────────────┐
│                    Customer Hardware                        │
│                                                             │
│  ┌─────────────────────── Management Plane ───────────────┐ │
│  │  Orchestrator                                          │ │
│  │  Persona Intelligence (soul)                           │ │
│  │  • Persona Store (versioned objects + memory)          │ │
│  │  • Behavior Engine                                     │ │
│  │  • Identity Manager                                    │ │
│  │  • Sensor & Trajectory Generator                       │ │
│  │  • Detection Radar                                     │ │
│  │  • License + API + Dashboard                           │ │
│  └────────────────────────────┬────────────────────────────┘ │
│                               │                             │
│                               ▼                             │
│  ┌─────────────────────── Execution Plane ────────────────┐ │
│  │  N × Isolated Redroid instances (bodies)               │ │
│  │  Each bound 1:1 to a living Persona                    │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

Optional Cloud (opt-in):
• Heavy DNA generation
• Model improvement
• Anonymized threat intelligence
```

**Invariant**: Persona Intelligence never lives inside the Redroid instance. The body is disposable. The soul is not.

---

## 2. Scaling Dimensions

We scale along three independent axes:

| Axis | Description | Year 1 Target | Later |
|------|-------------|---------------|-------|
| Personas per customer | Concurrent living personas | 5 → 50 | 200–1000+ |
| Customers | Independent self-hosted deployments | Dozens | Thousands |
| Capability depth | Presence → Conversation → MCP tools | Presence first | Full agentic |

Density (personas per GB RAM) is a constraint we improve, never a goal that sacrifices realism.

---

## 3. Key Scenarios — What Happens & What Is Needed

### 3.1 Birth of a New Persona

**Trigger**: Customer (or Natural Language Studio) requests a new persona.

**Flow**:
1. DNA Generator (local or optional cloud) produces a complete, coherent Persona object.
2. Persona is stored in the Persona Store (versioned, serializable).
3. Identity Manager creates a matching device profile.
4. Orchestrator schedules a Redroid instance.
5. Instance boots → Identity + Sensor hooks injected → Persona is bound.
6. Behavior Engine begins the “warm presence” phase (circadian-aware, low-risk activity).

**Needed**:
- Fast, high-quality DNA generation
- Persona object schema that is complete enough for long-term life
- Clean binding protocol between soul and body
- Resource reservation so the new body doesn’t starve existing ones

### 3.2 Daily Life of a Living Persona

**Every active persona continuously has**:
- A current high-level goal / schedule (from its circadian + memory + customer goals)
- Short-term context (what it just saw, did, felt)
- Long-term memory summaries

**Loop** (simplified):
1. Planner reads Persona state + current time + customer goals → emits next high-level intention.
2. Action Sequencer turns intention into a sequence of app-level actions.
3. Trajectory Generator produces realistic touch + sensor streams for those actions.
4. Streams are injected into the bound Redroid instance.
5. Screen state (and optional VLM) feeds back into short-term memory.
6. Significant events are summarized into longer-term memory.
7. Detection Radar watches for early warning signs and can trigger adaptation.

**Needed**:
- Hierarchical Behavior Engine that stays in character
- Efficient memory read/write (especially long-term summaries)
- Low-latency trajectory generation at scale
- Feedback path from body → soul without tight coupling

### 3.3 Scaling Number of Personas on One Machine

**What changes**:
- Orchestrator must pack instances intelligently (respecting circadian so not all are active at peak)
- Sensor/Trajectory generation must stay efficient (batching, shared models, or small on-device helpers)
- Persona Store and memory system must handle concurrent access cleanly
- Network / proxy isolation must not collapse

**What stays the same**:
- Each persona remains a fully coherent individual
- No shared brain that makes them start feeling coordinated
- Realism is not degraded to gain density

**Needed**:
- Strong resource accounting and admission control
- Circadian-aware scheduling
- Clear density vs realism measurement (as decided)

### 3.4 Multiple Customers (Horizontal Scale)

Each customer runs their own full stack on their own hardware.  
There is no central multi-tenant execution plane in Year 1.

**Optional central services** (opt-in only):
- DNA generation quality
- Model updates
- Anonymized Detection Radar intelligence

This keeps privacy, simplifies security, and matches the self-hosted sales motion.

### 3.5 Future: Conversation with a Persona

**When this arrives**:
- A conversation interface talks directly to the Persona Intelligence layer (not to the Redroid body).
- The same memory, personality, and history are used.
- The body may be idle or even powered down while the soul is in conversation.

**Needed later**:
- Clean conversational interface to the Persona object
- Ability to run “soul only” without a live body

### 3.6 Future: Persona Acts via MCP Tools

**When this arrives**:
- Persona Intelligence can emit tool calls through MCP servers.
- The same coherent identity and memory guide the tool use.
- Guardrails / policy module (separate) can constrain which tools or actions are allowed.

**Needed later**:
- Tool-calling interface on the Persona side
- Clear separation so tool use does not pollute the core presence engine

---

## 4. Critical Subsystems for Scale

### Persona Store
- Must support fast read of full Persona + selective memory retrieval
- Versioning and export (customer ownership)
- Concurrent access from Planner, Radar, Dashboard, future conversation layer

### Behavior Engine
- Must scale horizontally with number of active personas
- Hot path (trajectory) should be efficient; cold path (planning) can be heavier
- Must preserve distinctiveness — no accidental global style collapse at scale

### Orchestrator
- Admission control (will this new persona fit without hurting existing ones?)
- Circadian-aware placement and wake/sleep
- Health, restart, snapshot, bind/unbind soul ↔ body

### Memory System
- Hierarchical from day one (even if Year 1 implementation is simple)
- Designed so “years of life” remain possible later
- Summarization pipeline that keeps long-term memory bounded

### Detection Radar
- Per-persona and population-level views
- Must remain effective as numbers grow
- Local first; optional shared intelligence later

---

## 5. Scaling Principles (Non-Negotiable)

1. **Soul stays separate from body** — always.
2. **Realism is not sacrificed for density**.
3. **Each persona remains a distinct individual** even at population scale.
4. **Customer deployments are independent** (no forced multi-tenancy of execution).
5. **Future agentic capabilities plug in cleanly** without rewriting the presence core.
6. **Memory and identity are long-horizon by design**.

---

## 6. Phase-Aligned Scaling Path

**Phase 1 (MVP)**  
- 4–14 concurrent personas per machine  
- Simple scheduling  
- Basic memory  
- Manual or simple DNA  

**Phase 2**  
- Stronger Persona object + real memory  
- Much better trajectory & presence  
- Natural language DNA  
- Early perception testing  

**Phase 3+**  
- Higher density with same or better realism  
- Population-level coordination controls  
- Conversation interface  
- MCP tool use  
- Optional central intelligence improvements  

---

## 7. Open Technical Questions (to be turned into ADRs)

- Exact Persona schema and memory storage format
- Orchestrator language and internal event model
- How trajectory generation is parallelized across many personas
- How we measure and enforce the realism-vs-density curve in practice
- Binding protocol between Persona and Redroid instance (lifecycle states)

These will be decided and recorded as we implement.
