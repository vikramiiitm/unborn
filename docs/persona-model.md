# Persona Model (First-Class Object)

Every running instance is bound to a rich **Persona** object.  
This is the core product surface.

## Persona Structure (Conceptual)

```yaml
id: uuid
version: 1

# Identity & Demographics
demographics:
  age_range: [25, 28]
  gender_presentation: ...
  location: "Berlin, Germany"
  life_context: "Creative professional, lives alone, values slow living"

# Taste & Interests
interest_graph:
  primary: ["design", "analog photography", "slow living"]
  secondary: [...]
  aesthetic_style: "minimal, warm, film photography tones"

# Linguistic Style
language:
  primary_language: "en"
  dialect_region: "international English with light European flavor"
  emoji_habits: "sparse, intentional"
  vocabulary_register: "thoughtful, slightly poetic, avoids slang extremes"
  comment_philosophy: "meaningful over frequent"

# Temporal & Physical
circadian:
  typical_wake: "08:30"
  peak_activity: ["11:00-13:00", "19:00-22:00"]
  timezone: "Europe/Berlin"

physical_context_preferences:
  - hand_held_walking
  - desk_evening
  - cafe

# Engagement Style
engagement:
  type: "thoughtful_commenter"   # lurker | enthusiastic_sharer | quiet_reader | etc.
  risk_tolerance: low
  volume_preference: quality_over_quantity
  reply_latency_distribution: ...

# Memory & Evolution
memory:
  short_term: []          # recent interactions
  long_term_summaries: [] # evolving preferences, relationships, opinions
  last_updated: ...

# Visual Identity Signals (when profiles matter)
visual:
  profile_aesthetic: ...
  photo_style_hints: ...

# Device Binding
device_profile_id: ...    # linked coherent Android identity
```

## Design Principles

1. **Coherence first** — Every signal (what they watch, how long they stay, how they comment, when they are active, sensor micro-behavior) must feel like it comes from the same human.
2. **Long-term consistency** > short-term volume.
3. **Distinctiveness within niche** — Multiple personas from the same niche still feel like different people.
4. **Perception optimized** — Optimized for how real users interpret the behavior, not only for detection scores.
5. **Evolvable** — Personas age, change habits, and develop preferences over weeks/months.

## Generation Path

Natural language niche description  
→ Persona DNA Generator  
→ Full Persona object(s)  
→ Bound to Redroid instance(s) with matching device identity + behavior parameters.

Advanced customers can refine or author custom Persona DNA.
