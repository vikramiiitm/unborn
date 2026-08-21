package body

import (
	"fmt"
	"os"
)

// ADBHost is where Redroid ADB ports are reachable from the orchestrator process.
// Inside Docker Compose this is usually host.docker.internal (Linux: extra_hosts).
// Bare-metal local run: 127.0.0.1
func ADBHost() string {
	if h := os.Getenv("ADB_HOST"); h != "" {
		return h
	}
	return "127.0.0.1"
}

func ADBAddr(port int) string {
	return fmt.Sprintf("%s:%d", ADBHost(), port)
}
