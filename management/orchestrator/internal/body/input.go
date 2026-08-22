package body

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func adbShell(ctx context.Context, port int, args ...string) error {
	addr := ADBAddr(port)
	_ = exec.CommandContext(ctx, "adb", "connect", addr).Run()
	full := append([]string{"-s", addr, "shell"}, args...)
	out, err := exec.CommandContext(ctx, "adb", full...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("adb shell %v: %w (%s)", args, err, strings.TrimSpace(string(out)))
	}
	return nil
}

// Tap at pixel coordinates (device space).
func InputTap(ctx context.Context, port, x, y int) error {
	return adbShell(ctx, port, "input", "tap", fmt.Sprintf("%d", x), fmt.Sprintf("%d", y))
}

// Swipe from (x1,y1) to (x2,y2) over durationMs.
func InputSwipe(ctx context.Context, port, x1, y1, x2, y2, durationMs int) error {
	if durationMs <= 0 {
		durationMs = 300
	}
	return adbShell(ctx, port, "input", "swipe",
		fmt.Sprintf("%d", x1), fmt.Sprintf("%d", y1),
		fmt.Sprintf("%d", x2), fmt.Sprintf("%d", y2),
		fmt.Sprintf("%d", durationMs))
}

// KeyEvent sends an Android keycode (3=HOME, 4=BACK, 82=MENU, 26=POWER).
func InputKey(ctx context.Context, port int, keycode int) error {
	return adbShell(ctx, port, "input", "keyevent", fmt.Sprintf("%d", keycode))
}

// Text types ASCII text (spaces as %s).
func InputText(ctx context.Context, port int, text string) error {
	text = strings.ReplaceAll(text, " ", "%s")
	return adbShell(ctx, port, "input", "text", text)
}
