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
	addr := fmt.Sprintf("127.0.0.1:%d", adbPort)

	// Ensure connected
	_ = exec.CommandContext(ctx, "adb", "connect", addr).Run()

	// Wait briefly if device not ready
	deadline := time.Now().Add(15 * time.Second)
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
		return nil, fmt.Errorf("screencap: %w (is adb installed and device booted?)", err)
	}
	if len(out) < 100 || out[0] != 0x89 { // PNG magic
		// sometimes CRLF corruption; still return if looks big enough
		if len(out) < 1000 {
			return nil, fmt.Errorf("screencap returned empty or invalid data")
		}
	}
	return out, nil
}
