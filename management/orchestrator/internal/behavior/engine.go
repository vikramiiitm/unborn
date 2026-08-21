package behavior

import (
	"time"

	"github.com/vikramiiitm/unborn/management/orchestrator/internal/persona"
)

// Action is a high-level intention produced by the Behavior Engine.
type Action struct {
	Type        string            `json:"type"` // idle, scroll, view, like, comment, open_app
	App         string            `json:"app,omitempty"`
	DurationSec int               `json:"duration_sec,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Reason      string            `json:"reason,omitempty"`
}

// Engine turns a Persona + current time into the next action.
// Phase 1: simple rule-based / statistical. Later: hierarchical + AI.
type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

// NextAction returns a simple presence-oriented action based on persona engagement style and time.
func (e *Engine) NextAction(p *persona.Persona, now time.Time) Action {
	if p == nil {
		return Action{Type: "idle", DurationSec: 60, Reason: "no persona"}
	}

	// Very basic circadian gate: if outside rough active window, prefer idle.
	hour := now.Hour()
	active := hour >= 9 && hour <= 23
	if !active {
		return Action{
			Type:        "idle",
			DurationSec: 300,
			Reason:      "outside typical active hours",
		}
	}

	switch p.Engagement.Type {
	case persona.EngagementLurker:
		return Action{
			Type:        "scroll",
			App:         "feed",
			DurationSec: 90,
			Reason:      "lurker presence — observe more than act",
		}
	case persona.EngagementQuietReader:
		return Action{
			Type:        "view",
			App:         "content",
			DurationSec: 180,
			Reason:      "quiet reader — long dwell",
		}
	case persona.EngagementEnthusiasticSharer:
		return Action{
			Type:        "like",
			App:         "feed",
			DurationSec: 30,
			Reason:      "enthusiastic — light positive signal",
		}
	case persona.EngagementSelectiveEngager:
		return Action{
			Type:        "scroll",
			App:         "feed",
			DurationSec: 60,
			Reason:      "selective — browse before committing",
		}
	default: // thoughtful_commenter and fallback
		return Action{
			Type:        "view",
			App:         "content",
			DurationSec: 120,
			Reason:      "thoughtful — consume before possible comment",
		}
	}
}
