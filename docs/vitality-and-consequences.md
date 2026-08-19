# Vitality & Consequence System

## Philosophy

Every action a Persona takes should have consequences — just as they do for humans.

We do not want weightless, consequence-free agents.  
We want synthetic humans that live under real pressure: detection, perception, coherence, and survival.

This pressure is what forces the system to become genuinely good instead of just theatrically realistic.

---

## Vitality Score

Every living Persona carries a **Vitality Score** (0–100).

### What moves the score

| Category | Examples | Direction |
|----------|----------|-----------|
| Survival | Clean days, no CAPTCHA spikes, stable reach | ↑ |
| Coherence | Behavior stays in character, multi-modal consistency | ↑ |
| Perception | Human evaluators rate it as “feels real” | ↑ |
| Adaptation | Successfully recovered from a warning | ↑ (recovery) |
| Detection risk | Early Radar warnings, integrity signals | ↓ |
| Unnatural patterns | Timing, volume, or style that looks synthetic | ↓ |
| Repeated failures | Multiple serious incidents in short window | ↓↓ |
| Stagnation | Long periods with no meaningful coherent activity | mild ↓ |

Vitality is deliberately **slow-moving**. We avoid jittery life/death cycles.

---

## Consequence Levels

| Vitality Range | State | System Response |
|----------------|-------|------------------|
| 80–100 | Thriving | Normal full presence |
| 55–79 | Stable | Normal + light monitoring |
| 30–54 | Under Pressure | Trigger adaptation (behavior parameters, risk, schedule) |
| 10–29 | Critical | Pause body, keep soul, enter recovery protocol |
| 0–9 | Collapsed | Graceful shutdown of body. Soul is preserved and marked for review |

**Important rule:** We almost never destroy the Persona object itself.  
The soul (identity + memory + history) is preserved. Only the body is recycled.  
This protects customer ownership and long-term memory value.

---

## Design Principles

1. **Consequences are real** — actions accumulate.
2. **Adaptation before death** — the system first tries to help the Persona recover.
3. **Soul is precious** — memory and identity are not casually deleted.
4. **Modular** — Vitality is a separate module that reads from Detection Radar, Behavior Engine, and (later) perception scores. It does not live inside the core Persona schema as punishment logic.
5. **Customer-configurable aggressiveness** — some customers will want very strict quality filters; others will prefer maximum presence.
6. **Transparent** — customers can see why Vitality changed.

---

## Future Evolution

Later, Personas themselves can become aware of their Vitality and attempt self-correction (especially once conversational + MCP capabilities exist).  
For Year 1 we keep it as an external system pressure.
