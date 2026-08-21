#!/usr/bin/env bash
# Unborn host installer — Ubuntu-focused Phase 1
set -euo pipefail

echo "=== Unborn installer ==="

if [[ "$(id -u)" -ne 0 ]]; then
  echo "Re-run with sudo for kernel modules and Docker setup."
  exit 1
fi

export DEBIAN_FRONTEND=noninteractive

echo "[1/5] Packages..."
apt-get update -qq
apt-get install -y -qq ca-certificates curl git linux-modules-extra-$(uname -r) || true

echo "[2/5] Kernel modules (binder / ashmem)..."
modprobe binder_linux devices="binder,hwbinder,vndbinder" 2>/dev/null || echo "  warn: binder_linux not loaded (check kernel)"
modprobe ashmem_linux 2>/dev/null || echo "  warn: ashmem_linux not loaded (ok on some kernels)"

if ! grep -q binder_linux /etc/modules 2>/dev/null; then
  echo "binder_linux devices=binder,hwbinder,vndbinder" >> /etc/modules || true
fi

# Modern Linux kernels (Kernel 6.x / 7.x on Ubuntu 24.04+) require mounting binderfs
mkdir -p /dev/binderfs
mount -t binder binder /dev/binderfs 2>/dev/null || true
ln -sf /dev/binderfs/binder /dev/binder 2>/dev/null || true
ln -sf /dev/binderfs/hwbinder /dev/hwbinder 2>/dev/null || true
ln -sf /dev/binderfs/vndbinder /dev/vndbinder 2>/dev/null || true

if ! grep -q "/dev/binderfs" /etc/fstab 2>/dev/null; then
  echo "binder /dev/binderfs binder defaults 0 0" >> /etc/fstab || true
fi

echo "[3/5] Docker..."
if ! command -v docker >/dev/null 2>&1; then
  curl -fsSL https://get.docker.com | sh
fi
systemctl enable --now docker 2>/dev/null || true

echo "[4/5] Docker Compose plugin..."
if ! docker compose version >/dev/null 2>&1; then
  apt-get install -y -qq docker-compose-plugin 2>/dev/null || true
fi

echo "[5/5] Data directories..."
mkdir -p /var/lib/unborn/redroid-data
chmod 755 /var/lib/unborn /var/lib/unborn/redroid-data

echo ""
echo "=== Host prep done ==="
echo "Next:"
echo "  cd docker && docker compose up --build -d"
echo "  open http://localhost:8080/"
echo ""
echo "Real Redroid bodies need: privileged containers + binder modules (above)."
echo "Start with simulated bodies first: POST /v1/instances {\"simulated\": true}"
