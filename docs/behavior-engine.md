# Behavior Engine

The Behavior Engine turns a Persona into concrete actions on a Redroid instance.

## Hierarchical Design

```
Persona (goals, personality, circadian, memory)
        ↓
High-level Planner
  - Daily/weekly schedule
  - Current goals (“warm account”, “engage with design content”, …)
  - Context awareness (time of day, recent history)
        ↓
Mid-level Action Sequencer
  - Chooses sequences of high-level actions (open app, scroll feed, watch video, leave comment, …)
  - Respects engagement philosophy and risk tolerance
        ↓
Low-level Trajectory Generator
  - Realistic touch paths, velocities, pressures
  - Correlated sensor streams
  - Micro-timing jitter that matches the persona’s physical context
```

## Personality → Action Mapping

Different personas in the same niche must still feel distinct:
- One is a quiet lurker who rarely comments but stays long on good content
- Another leaves thoughtful, slightly longer comments
- Another shares more enthusiastically but only within strict taste boundaries

The engine samples from the persona’s distributions rather than running generic scripts.

## Future AI Layers
- Small vision models for actual screen understanding
- Lightweight planners / agents that can replan when the UI changes or friction appears
- On-device tiny models for low-latency next-action decisions

## Consistency Rules
- All modalities (visual attention, dwell time, likes, comments, sensor behavior, activity times) must stay coherent with the same Persona object.
- Long-term memory influences future preferences and avoidance patterns.
