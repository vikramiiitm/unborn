package orchestrator

import (
	"context"

	"github.com/vikramiiitm/unborn/management/orchestrator/internal/body"
)

func (o *Orchestrator) InstallAPK(ctx context.Context, bodyID string, apkPath string) (string, error) {
	port, err := o.realADBPort(bodyID)
	if err != nil {
		return "", err
	}
	return body.InstallAPK(ctx, port, apkPath)
}

func (o *Orchestrator) InstallAPKBytes(ctx context.Context, bodyID string, data []byte, name string) (string, error) {
	port, err := o.realADBPort(bodyID)
	if err != nil {
		return "", err
	}
	return body.InstallAPKBytes(ctx, port, data, name)
}
