package body

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ScreenshotPNG captures the device screen via adb screencap (real bodies only).
func ScreenshotPNG(ctx context.Context, adbPort int) ([]byte, error) {
	if adbPort <= 0 {
		return nil, fmt.Errorf("no adb port")
	}
	addr := ADBAddr(adbPort)

	_ = exec.CommandContext(ctx, "adb", "connect", addr).Run()

	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		out, _ := exec.CommandContext(ctx, "adb", "-s", addr, "get-state").CombinedOutput()
		if strings.Contains(string(out), "device") {
			break
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(1 * time.Second):
		}
	}

	cmd := exec.CommandContext(ctx, "adb", "-s", addr, "exec-out", "screencap", "-p")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("screencap via %s: %w (adb in image? body booted? ADB_HOST set?)", addr, err)
	}
	if len(out) < 1000 {
		return nil, fmt.Errorf("screencap returned empty or invalid data from %s", addr)
	}
	return out, nil
}
