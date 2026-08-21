package orchestrator

import (
	"context"
	"log"
	"time"

	"github.com/vikramiiitm/unborn/management/orchestrator/internal/body"
)

// RunBehaviorLoop periodically picks next actions for personas that have a running body.
// Phase 1: log + light vitality nudge. Later: inject into Redroid / drive UI.
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
	activePersona := map[string]bool{}
	for _, inst := range instances {
		if inst.State == body.StateRunning {
			activePersona[inst.PersonaID] = true
		}
	}
	if len(activePersona) == 0 {
		return
	}

	for personaID := range activePersona {
		act, err := o.NextAction(ctx, personaID)
		if err != nil {
			log.Printf("behavior: persona %s: %v", personaID, err)
			continue
		}
		log.Printf("behavior: persona=%s action=%s app=%s duration=%ds reason=%s",
			personaID, act.Type, act.App, act.DurationSec, act.Reason)

		// Light presence reward — real Radar will drive serious vitality changes later
		o.vitality.Adjust(personaID, 0.05, "presence_tick:"+act.Type)
	}
}
