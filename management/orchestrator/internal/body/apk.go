package body

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// InstallAPK pushes and installs an APK on a real body via adb.
// apkPath must be a readable path on the orchestrator host/container.
func InstallAPK(ctx context.Context, adbPort int, apkPath string) (string, error) {
	if adbPort <= 0 {
		return "", fmt.Errorf("no adb port")
	}
	if apkPath == "" {
		return "", fmt.Errorf("apk path required")
	}
	st, err := os.Stat(apkPath)
	if err != nil || st.IsDir() {
		return "", fmt.Errorf("apk not found: %s", apkPath)
	}
	addr := ADBAddr(adbPort)
	_ = exec.CommandContext(ctx, "adb", "connect", addr).Run()

	// adb install -r -g (replace + grant runtime permissions where possible)
	cmd := exec.CommandContext(ctx, "adb", "-s", addr, "install", "-r", "-g", apkPath)
	out, err := cmd.CombinedOutput()
	s := strings.TrimSpace(string(out))
	if err != nil {
		return s, fmt.Errorf("adb install: %w (%s)", err, s)
	}
	if !strings.Contains(strings.ToLower(s), "success") {
		return s, fmt.Errorf("install may have failed: %s", s)
	}
	return s, nil
}

// InstallAPKBytes writes bytes to a temp file then installs.
func InstallAPKBytes(ctx context.Context, adbPort int, data []byte, name string) (string, error) {
	if len(data) < 100 {
		return "", fmt.Errorf("apk data too small")
	}
	if name == "" {
		name = "app.apk"
	}
	dir, err := os.MkdirTemp("", "unborn-apk-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, filepath.Base(name))
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return InstallAPK(ctx, adbPort, path)
}
