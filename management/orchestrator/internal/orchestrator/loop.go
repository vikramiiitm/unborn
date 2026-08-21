package orchestrator

import (
	"context"
	"log"
	"time"

	"github.com/vikramiiitm/unborn/management/orchestrator/internal/body"
)

func (o *Orchestrator) RunBehaviorLoop(ctx context.Context) {
	if !o.cfg.BehaviorLoop {
		log.Println("Behavior loop disabled")
		return
	}
	interval := time.Duration(o.cfg.BehaviorIntervalSec) * time.Second
	log.Printf("Behavior loop started (every %s)", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Behavior loop stopped")
			return
		case <-ticker.C:
		o.tickBehavior(ctx)
		}
	}
}

func (o *Orchestrator) tickBehavior(ctx context.Context) {
	instances := o.bodies.List()
	for _, inst := range instances {
		if inst.State != body.StateRunning {
			continue
		}
		personaID := inst.PersonaID

		// Occasional health check on real bodies
		if !inst.Simulated && inst.ADBPort > 0 {
			ok, reason := o.CheckBodyHealth(ctx, inst.ID)
			if !ok {
				log.Printf("health: body=%s persona=%s unhealthy: %s", inst.ID[:8], personaID, reason)
				o.vitality.Adjust(personaID, -0.5, "adb_unhealthy")
			}
		}

		act, err := o.NextAction(ctx, personaID)
		if err != nil {
			log.Printf("behavior: persona %s: %v", personaID, err)
			continue
		}
		log.Printf("behavior: persona=%s action=%s app=%s duration=%ds reason=%s",
			personaID, act.Type, act.App, act.DurationSec, act.Reason)
		o.vitality.Adjust(personaID, 0.05, "presence_tick:"+act.Type)
	}
}
