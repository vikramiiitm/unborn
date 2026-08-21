# Redroid Runtime (Execution Plane)

Unborn uses [Redroid](https://github.com/remote-android/redroid-doc) as the Android **body**. Personas (souls) live in the Management Plane; each running body is a privileged Docker container.

---

## Host requirements (Ubuntu)

```bash
sudo apt install linux-modules-extra-$(uname -r)
sudo modprobe binder_linux devices="binder,hwbinder,vndbinder"
sudo modprobe ashmem_linux   # if available on your kernel

# persist
echo 'binder_linux devices="binder,hwbinder,vndbinder"' | sudo tee -a /etc/modules
```

Docker must allow privileged containers. Orchestrator needs access to the Docker socket (see `docker/docker-compose.yml`).

Optional but recommended:
```bash
sudo apt install android-tools-adb   # for health checks
```

---

## What Unborn starts

For each **real** instance:

```text
docker run -d --privileged
  --name unborn-<id>
  --memory / --cpus limits (configurable)
  -v <dataRoot>/<bodyId>:/data
  -p <hostPort>:5555
  redroid/redroid:14.0.0_64only-latest
  androidboot.redroid_width=1080
  androidboot.redroid_height=1920
  androidboot.redroid_dpi=480
  [optional proxy boot props]
```

- **ADB** exposed on host ports starting at `5555`
- **Data** persisted under `/var/lib/unborn/redroid-data/<bodyId>` (or compose volume)
- **Proxy** from persona assignment → `androidboot.redroid_net_proxy_*`

Simulated bodies skip Docker entirely (control-plane dev).

---

## Images

Default: `redroid/redroid:14.0.0_64only-latest`

Also common: `12.0.0_64only-latest`, `15.0.0_64only-latest`, `16.0.0-latest` — set via env `REDROID_IMAGE` later.

---

## Lifecycle states

`pending → starting → running → stopping → stopped | failed`

API:
- `POST /v1/instances` `{ "persona_id", "simulated": false }` + `?real=true`
- `POST /v1/instances/{id}/stop`
- `GET /v1/instances/{id}/health` — ADB reachability

---

## Density (starting guidelines)

| Host RAM | Concurrent real Redroid (rough) |
|----------|----------------------------------|
| 16 GB | 4–6 |
| 32 GB | 10–14 |
| 64 GB+ | 20+ (tune memory limits) |

Always pair with Persona coherence — density without souls is a normal farm.

---

## Next Redroid work

1. Per-body memory/CPU limits from config
2. Device profile → build props / identity injection
3. Frida / sensor injection hooks
4. Network namespace isolation beyond proxy
5. Snapshot / wipe data volume API
