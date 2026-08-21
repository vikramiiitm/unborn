# Unborn

**Farming with souls.**

Self-hosted platform: living populations of digital people on your Ubuntu hardware.  
Every body (Redroid) is driven by a Persona — taste, memory, circadian rhythm, engagement style, and **Vitality** (real consequences).

> Other farms rent you activity.  
> Unborn runs a living population of people who generate activity.

---

## Quick start (dev)

### Requirements
- Ubuntu 22.04+ recommended (or similar)
- Docker + Docker Compose plugin
- Optional for **real** Redroid bodies: kernel modules + `adb`

### 1. Host prep (once, for real Android bodies)

```bash
sudo bash installer/install.sh
# loads binder/ashmem, installs Docker if needed, creates data dirs
```

Or manually:

```bash
sudo apt install linux-modules-extra-$(uname -r) android-tools-adb
sudo modprobe binder_linux devices="binder,hwbinder,vndbinder"
sudo modprobe ashmem_linux   # if available
```

### 2. Start the stack

```bash
git clone https://github.com/vikramiiitm/unborn.git
cd unborn/docker
docker compose up --build -d
```

Services:
| Service | Port |
|---------|------|
| Orchestrator + Dashboard | **8080** |
| Postgres | 5432 |
| Redis | 6379 |

### 3. Open the UI

```text
http://localhost:8080/
```

Tabs: **Population** · **Bodies** · **Playbooks** · **Proxies**  
Header shows license tier and max instances.

### 4. API examples

```bash
# Health
curl http://localhost:8080/health

# Create persona
curl -X POST http://localhost:8080/v1/personas \
  -H 'Content-Type: application/json' \
  -d '{"display_name":"Ava","location":"Berlin","timezone":"Europe/Berlin","engagement":"thoughtful_commenter"}'

# Simulated body (no Redroid needed)
curl -X POST http://localhost:8080/v1/instances \
  -H 'Content-Type: application/json' \
  -d '{"persona_id":"<PERSONA_ID>","simulated":true}'

# Real Redroid body (needs modules + docker.sock)
curl -X POST 'http://localhost:8080/v1/instances?real=true' \
  -H 'Content-Type: application/json' \
  -d '{"persona_id":"<PERSONA_ID>","simulated":false}'

# Optional: proxy before real start
curl -X PUT http://localhost:8080/v1/personas/<PERSONA_ID>/proxy \
  -H 'Content-Type: application/json' \
  -d '{"host":"1.2.3.4","port":1080,"type":"http"}'

# License / vitality / playbooks
curl http://localhost:8080/v1/license
curl http://localhost:8080/v1/vitality
curl http://localhost:8080/v1/playbooks
```

### 5. Stop

```bash
cd docker && docker compose down
```

---

## Real Redroid notes

- Orchestrator mounts `/var/run/docker.sock` to start privileged Redroid containers.
- Default image: `redroid/redroid:14.0.0_64only-latest`
- Env overrides: `REDROID_IMAGE`, `REDROID_MEMORY`, `REDROID_CPUS`, `REDROID_DATA_ROOT`, `REDROID_BASE_ADB_PORT`
- After boot: identity is pushed over ADB (model, manufacturer, android_id, …)
- Docs: **[docs/redroid.md](docs/redroid.md)**

**Start with simulated bodies** until the host modules and image pull are verified.

---

## Status & docs

| File | Purpose |
|------|---------|
| [STATUS.md](STATUS.md) | What is built right now |
| [FEATURES.md](FEATURES.md) | Feature checklist |
| [docs/README.md](docs/README.md) | Doc index |
| [docs/farming-with-souls.md](docs/farming-with-souls.md) | Positioning |
| [docs/management-ux.md](docs/management-ux.md) | Dashboard / automations vision |
| [docs/redroid.md](docs/redroid.md) | Execution plane |

---

## Repo layout

```
unborn/
├── docker/                 # Compose (Postgres, Redis, Orchestrator)
├── installer/install.sh    # Host prep
├── management/orchestrator # Go control plane + dashboard
├── docs/                   # Strategy + architecture
├── STATUS.md / FEATURES.md
└── README.md
```

---

*Unborn — digital people, born on your hardware.*
