# Redroid Runtime (Execution Plane)

Unborn uses [Redroid](https://github.com/remote-android/redroid-doc) as the Android **body**. Personas live in the Management Plane; each running body is a privileged Docker container.

See also: root **[README.md](../README.md)** (how to run), **[management-ux.md](management-ux.md)** (dashboard).

---

## Host requirements (Ubuntu)

```bash
sudo bash installer/install.sh
# or manually:
sudo apt install linux-modules-extra-$(uname -r) android-tools-adb
sudo modprobe binder_linux devices="binder,hwbinder,vndbinder"
sudo modprobe ashmem_linux   # if available
```

Orchestrator needs Docker socket (mounted in `docker/docker-compose.yml`).

---

## What gets started (real body)

```text
docker run -d --privileged
  --name unborn-<shortId>
  --memory 3072m --cpus 2.0   # env-overridable
  -v <dataRoot>/<bodyId>:/data
  -p <hostPort>:5555
  redroid/redroid:14.0.0_64only-latest
  androidboot.redroid_width=1080
  androidboot.redroid_height=1920
  androidboot.redroid_dpi=480
  [androidboot.redroid_net_proxy_* if persona has proxy]
```

After boot (async):
- ADB identity inject: `android_id`, `ro.product.model`, manufacturer, brand, serial

Simulated bodies skip Docker (dev / control-plane only).

---

## Environment

| Env | Default | Meaning |
|-----|---------|--------|
| `REDROID_IMAGE` | `redroid/redroid:14.0.0_64only-latest` | Image tag |
| `REDROID_DATA_ROOT` | `/var/lib/unborn/redroid-data` | Host data dirs |
| `REDROID_MEMORY` | `3072m` | Container memory limit |
| `REDROID_CPUS` | `2.0` | CPU limit |
| `REDROID_BASE_ADB_PORT` | `5555` | First host ADB port |

---

## Instance API (bodies)

| Method | Path | Notes |
|--------|------|--------|
| POST | `/v1/instances` | `{ "persona_id", "simulated" }` — use `?real=true` for Redroid |
| GET | `/v1/instances` | List |
| GET | `/v1/instances/{id}` | Detail (adb_port, container_name, data_dir, …) |
| POST | `/v1/instances/{id}/stop` | Stop / remove container |
| POST | `/v1/instances/{id}/restart` | `docker restart` (real only) |
| POST | `/v1/instances/{id}/wipe` | Delete data dir (must be stopped) |
| GET | `/v1/instances/{id}/health` | Container running + ADB |
| GET | `/v1/instances/{id}/logs?tail=100` | Docker logs (text) |
| POST | `/v1/instances/{id}/inject-identity` | Re-run ADB identity inject |

Proxy for next real start: `PUT /v1/personas/{id}/proxy`.

---

## Density (guidelines)

| Host RAM | Concurrent real Redroid (approx.) |
|----------|-----------------------------------|
| 16 GB | 4–6 |
| 32 GB | 10–14 |
| 64 GB+ | 20+ (tune limits) |

Density without Personas is a normal farm — coherence stays the product.

---

## Still open

1. Stronger fingerprint / persistent build.prop strategy
2. Network namespace isolation beyond proxy boot props
3. Frida / touch / sensor injection
4. Snapshot API
