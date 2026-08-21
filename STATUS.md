# Unborn – Project Status

Last updated: 2026-08-22

## Focus: Redroid (Execution Plane)

### Redroid now
- Real `docker run` with memory/CPU limits, data dir, ADB port, proxy boot props
- Env: `REDROID_IMAGE`, `REDROID_DATA_ROOT`, `REDROID_MEMORY`, `REDROID_CPUS`, `REDROID_BASE_ADB_PORT`
- Health: container inspect + ADB
- `POST /v1/instances/{id}/wipe` after stop
- Docs: `docs/redroid.md`

### Still on Redroid track
1. Real device identity injection (build.prop / settings via adb)
2. Network namespace isolation
3. Frida / input injection path
4. Snapshot API

Personas stay as-is until Redroid path is solid.
