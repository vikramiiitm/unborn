# Pre-Build Decisions (Locked)

These answers are now the default for Unborn. We can revisit later only with strong evidence.

---

## 1. Cloud vs Self-hosted Split for Persona Intelligence

**Decision: Hybrid with strong customer-side core**

- **Runs fully on customer hardware (self-hosted):**
  - Persona runtime & memory store
  - Behavior Engine (planner + sequencer + trajectory)
  - Identity Manager
  - Sensor / touch generation
  - Detection Radar (local signals)
  - Orchestrator & lifecycle

- **Can use optional cloud / central services (customer opt-in):**
  - Heavy Persona DNA generation from natural language (larger models)
  - Periodic model improvements / fine-tunes
  - Anonymized threat intelligence sharing
  - Advanced perception scoring helpers

**Why this split:**
- Privacy and “runs on my metal” is a strong selling point for agencies and brands.
- We keep the most sensitive and latency-critical pieces local.
- We still get a data flywheel and superior DNA quality through optional central intelligence.
- Matches the “engine lives separately” principle cleanly.

Default posture: everything works offline/air-gapped. Cloud is pure upside.

---

## 2. How We Measure “Feels Real”

**Decision: Dual metric system from day one**

1. **Technical / Survival metrics** (necessary but not sufficient)
   - Account longevity
   - CAPTCHA / integrity / reach signals
   - Detection Radar alerts

2. **Human Perception metrics** (the real north star)
   - Blind human evaluation: “Does this feel like a real person from this niche?”
   - Structured scoring on coherence, presence, distinctiveness, comment quality, timing
   - Small regular panels (even 5–10 people) reviewing short interaction traces
   - Later: customer-side perception feedback loops

**Rule:**  
We never declare a behavior or persona “good” based only on surviving detection.  
If humans can tell, it is not good enough — even if the detectors are currently fooled.

We will build lightweight perception testing into the development loop early (Phase 2).

---

## 3. Hard Ethical Boundaries

**Decision: Explicit refusal list + positive principles**

### We refuse:
- Personas designed to impersonate real specific private individuals
- Child / minor personas (anyone under 18)
- Personas whose primary purpose is harassment, scams, or non-consensual deep involvement in real people’s lives
- Political influence operations or anything that looks like coordinated inauthentic behavior for electoral purposes
- Anything that requires breaking the law in the customer’s jurisdiction

### Positive principles:
- Transparency with the customer about what the system is
- Strong preference for use cases that add genuine social proof or research value rather than pure deception at scale
- Customer remains responsible for how they deploy personas
- We will evolve this list as the technology and regulation change, but the core refusals stay.

These boundaries protect long-term legitimacy and make enterprise / serious brand sales much easier.

---

## 4. Persona DNA Ownership & IP Model

**Decision: Customer owns their custom Personas. We own the engine and improvements.**

- Any Persona DNA a customer creates or heavily customizes belongs to the customer.
- They can export it, delete it, or take it with them.
- We may use **anonymized, aggregated** learnings to improve the core models (opt-out available).
- Base personas / templates and the generation system remain ours.
- Clear contractual language for agencies managing personas on behalf of their clients.

This is the cleanest model for trust and for closing agency + brand deals.

---

## 5. Density vs Realism Trade-off

**Decision: Realism wins. Density is a constraint we improve, not a goal we sacrifice for.**

- We will never ship a mode that makes personas feel generic or coordinated just to fit more instances on a box.
- Published density numbers are always paired with the realism level they were measured at.
- If forced to choose in the short term, we choose fewer higher-quality personas over many mediocre ones.
- Long-term engineering goal is to raise density **while holding or improving** coherence and presence.

This protects the core moat. Most competitors will make the opposite choice. That is fine.

---

## Summary of Posture

- Engine (Persona Intelligence) lives separately and is the valuable part.
- Customer hardware runs the living personas by default.
- Optional cloud makes the DNA and models better.
- We optimize for human perception of realism + survival, not just detector scores.
- Clear ethical lines.
- Customer owns their people.
- Quality > cheap density.

These decisions keep Year 1 focused and commercially sharp while leaving the door wide open for the much larger future.
