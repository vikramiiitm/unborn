# Architecture

## High-Level Planes

### Management Plane
Central control system running as Docker services on the Ubuntu host.

- **Orchestrator**  
  Lifecycle management of personas + Redroid instances.  
  Schedules, health checks, resource allocation, snapshots, resets.

- **API** (REST + gRPC)  
  Programmatic control for customers and internal agents.

- **Web Dashboard**  
  Human operators create/monitor personas, view Detection Radar, manage licenses.

- **License Service**  
  Offline-capable signed keys (JWT / commercial system). Enforces max instances, tiers, expiry.

- **Identity Manager**  
  Generates and injects coherent device fingerprints that match the persona’s claimed device model, sensors, GPU, installed apps baseline, etc.

- **Sensor & Input Provider**  
  Generates correlated multi-modal streams (accelerometer, gyroscope, touch trajectories, ambient light, etc.) driven by persona physical context (hand-held, desk, walking, commuting…).

- **Behavior Engine**  
  Hierarchical:  
  1. High-level Planner (goals + daily schedule from persona)  
  2. Mid-level Action Sequencer  
  3. Low-level Trajectory Generator (realistic touch + sensor micro-behavior)

- **Learning & Adaptation Layer**  
  Improves personas over time from reception signals, detection events, and synthetic data.

- **Observability & Detection Radar**  
  Early warning signs (CAPTCHA spikes, reach drops, integrity failures). Strategy switching / backoff.

### Execution Plane
N isolated Redroid instances.

Each instance:
- Bound to exactly one Persona
- Unique network namespace + forced proxy (mobile/residential preferred)
- Unique full device identity
- Sensor/touch injection hooks (Frida → later custom HAL)
- Optional tiny on-device models for low-latency decisions

## Network Isolation
- Per-instance Docker network or network namespace
- Forced routing through customer-supplied proxy
- Redroid native proxy + Frida-level leak protection (WebRTC, DNS, STUN)
- Sticky sessions + rotation support
- Different exit nodes per persona when needed

## Resource Guidelines (starting point)
- 16 GB RAM → ~4–6 concurrent instances
- 32 GB RAM → ~10–14 concurrent instances
- Strict Docker resource limits + health checks

## Delivery Model
- One-command installer
- Docker Compose based
- License-key activated
- Self-hosted on customer Ubuntu hardware
- Easy update path
