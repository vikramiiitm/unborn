package orchestrator

import (
	"context"

	"github.com/vikramiiitm/unborn/management/orchestrator/internal/body"
)

func (o *Orchestrator) realADBPort(bodyID string) (int, error) {
	b, ok := o.bodies.Get(bodyID)
	if !ok {
		return 0, ErrInstanceNotFound
	}
	if b.Simulated {
		return 0, ErrSimulatedNoScreen
	}
	if b.ADBPort <= 0 {
		return 0, ErrInstanceNotFound
	}
	return b.ADBPort, nil
}

func (o *Orchestrator) InputTap(ctx context.Context, bodyID string, x, y int) error {
	port, err := o.realADBPort(bodyID)
	if err != nil {
		return err
	}
	return body.InputTap(ctx, port, x, y)
}

func (o *Orchestrator) InputSwipe(ctx context.Context, bodyID string, x1, y1, x2, y2, ms int) error {
	port, err := o.realADBPort(bodyID)
	if err != nil {
		return err
	}
	return body.InputSwipe(ctx, port, x1, y1, x2, y2, ms)
}

func (o *Orchestrator) InputKey(ctx context.Context, bodyID string, keycode int) error {
	port, err := o.realADBPort(bodyID)
	if err != nil {
		return err
	}
	return body.InputKey(ctx, port, keycode)
}

func (o *Orchestrator) InputText(ctx context.Context, bodyID string, text string) error {
	port, err := o.realADBPort(bodyID)
	if err != nil {
		return err
	}
	return body.InputText(ctx, port, text)
}
