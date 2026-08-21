# Troubleshooting & Kernel Diagnostics

This guide documents common operational issues, Linux kernel setup requirements for Redroid, port conflict resolutions, and diagnostic steps for the **Unborn** platform.

---

## 1. Redroid Kernel Panic & Missing `/dev/binder` Nodes

### Symptom
- PC kernel panics/freezes when running real Redroid containers (`redroid/redroid:14.0.0_64only-latest`).
- Redroid container crashes internally on startup (`servicemanager` or `vold` stuck in `[crash_dump64]`).
- Error in `ls`: `ls: cannot access '/dev/binder': No such file or directory`.

### Cause
On modern Linux kernels (**Kernel 6.x and 7.x on Ubuntu 24.04+**), loading `binder_linux` registers the `binder` filesystem in kernel memory, but does not automatically create raw device nodes in `/dev/`. Redroid requires active Android IPC devices (`/dev/binder`, `/dev/hwbinder`, `/dev/vndbinder`).

### Fix
Execute the following commands to mount **binderfs** and create device symlinks:

```bash
# 1. Load binder kernel module
sudo modprobe binder_linux devices="binder,hwbinder,vndbinder"

# 2. Mount binderfs
sudo mkdir -p /dev/binderfs
sudo mount -t binder binder /dev/binderfs

# 3. Create symlinks for Redroid container access
sudo ln -sf /dev/binderfs/binder /dev/binder
sudo ln -sf /dev/binderfs/hwbinder /dev/hwbinder
sudo ln -sf /dev/binderfs/vndbinder /dev/vndbinder
```

#### Make it permanent across host reboots:
Add `binder_linux` to `/etc/modules` and `binderfs` to `/etc/fstab`:
```bash
echo "binder_linux devices=binder,hwbinder,vndbinder" | sudo tee -a /etc/modules
echo "binder /dev/binderfs binder defaults 0 0" | sudo tee -a /etc/fstab
```

#### Verification
```bash
ls -l /dev/binder /dev/hwbinder /dev/vndbinder
# Expected output:
# /dev/binder -> /dev/binderfs/binder
# /dev/hwbinder -> /dev/binderfs/hwbinder
# /dev/vndbinder -> /dev/binderfs/vndbinder
```

---

## 2. Docker Port Allocation Failure (`Bind for 0.0.0.0:5555 failed`)

### Symptom
```text
Fail: docker run failed: exit status 125 (...): Bind for 0.0.0.0:5555 failed: port is already allocated.
```

### Cause
- Host port `5555` is occupied by a host `adb` server daemon, another local container, or a lingering Redroid instance.
- Orchestrator port allocation originally checked only running containers in `docker ps`, missing ports occupied by stopped containers or host TCP listeners.

### Resolution
1. **Automated Fix**: The Orchestrator `allocatePort()` logic in `management/orchestrator/internal/body/redroid.go` now checks `docker ps -a` and performs an active TCP socket probe (`net.Listen`) before picking a port, automatically skipping busy host ports (e.g. assigning `5556` if `5555` is in use).
2. **Manual Cleanup of Stuck Containers**:
   ```bash
   docker ps -a | grep redroid
   docker rm -f <container_name_or_id>
   ```

---

## 3. Dashboard Displaying "Waiting for boot / adb…"

### Symptom
Real body thumbnail in Dashboard shows `<div class="ph fail">Waiting for boot / adb…</div>`.

### Diagnostic Steps
1. **Check Container Status**:
   ```bash
   docker ps -a
   ```
2. **Inspect Internal Processes**:
   ```bash
   docker exec <container_name> ps -ef
   ```
   - If processes like `[servicemanager]`, `[hwservicemanage]`, or `[crash_dump64]` are present, **BinderFS is missing or unmounted on the host** (see Section 1).
3. **Test ADB Connectivity**:
   ```bash
   docker exec unborn-orchestrator adb connect host.docker.internal:<adbPort>
   ```
   - If `Connection refused`, Redroid's internal `adbd` daemon did not start (due to missing binder driver).

---

## 4. Simulated vs Real Bodies Usage Guidelines

| Mode | Memory / Overhead | Use Case | Host Requirements |
| :--- | :--- | :--- | :--- |
| **Simulated** (`simulated: true`) | ~0.01% RAM | High-density persona testing (100+ instances), CI/CD pipelines, API/behavior logic dev | Any OS / Docker host |
| **Real** (`real=true`) | ~3 GB RAM / body | Full Android OS execution, ADB screen capture, UI automation, app interactions | Linux + `binderfs` kernel drivers |
