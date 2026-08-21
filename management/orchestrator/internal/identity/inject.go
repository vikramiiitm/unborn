package identity

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/vikramiiitm/unborn/management/orchestrator/internal/body"
)

func InjectViaADB(ctx context.Context, adbPort int, p *DeviceProfile) error {
	if p == nil || adbPort <= 0 {
		return fmt.Errorf("invalid profile or port")
	}
	addr := body.ADBAddr(adbPort)

	deadline := time.Now().Add(45 * time.Second)
	for time.Now().Before(deadline) {
		_ = exec.CommandContext(ctx, "adb", "connect", addr).Run()
		out, _ := exec.CommandContext(ctx, "adb", "-s", addr, "get-state").CombinedOutput()
		if strings.Contains(string(out), "device") {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}

	cmds := [][]string{
		{"shell", "settings", "put", "secure", "android_id", sanitize(p.AndroidID)},
		{"shell", "setprop", "ro.product.model", sanitize(p.Model)},
		{"shell", "setprop", "ro.product.manufacturer", sanitize(p.Manufacturer)},
		{"shell", "setprop", "ro.product.brand", sanitize(p.Manufacturer)},
	}
	if p.Serial != "" {
		cmds = append(cmds, []string{"shell", "setprop", "ro.serialno", sanitize(p.Serial)})
	}
	for k, v := range p.Properties {
		if k != "" && v != "" {
			cmds = append(cmds, []string{"shell", "setprop", k, sanitize(v)})
		}
	}

	var lastErr error
	for _, args := range cmds {
		full := append([]string{"-s", addr}, args...)
		out, err := exec.CommandContext(ctx, "adb", full...).CombinedOutput()
		if err != nil {
			lastErr = fmt.Errorf("adb %v: %w (%s)", args, err, strings.TrimSpace(string(out)))
		}
	}
	return lastErr
}

func sanitize(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "_")
	return s
}
