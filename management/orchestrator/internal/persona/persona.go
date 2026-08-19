package persona

import (
	"time"

	"github.com/google/uuid"
)

// Status represents the lifecycle of a Persona soul.
type Status string

const (
	StatusDraft    Status = "draft"
	StatusWarming  Status = "warming"
	StatusActive   Status = "active"
	StatusPaused   Status = "paused"
	StatusArchived Status = "archived"
)

// EngagementType defines how the persona tends to interact.
type EngagementType string

const (
	EngagementLurker             EngagementType = "lurker"
	EngagementThoughtfulCommenter EngagementType = "thoughtful_commenter"
	EngagementEnthusiasticSharer EngagementType = "enthusiastic_sharer"
	EngagementQuietReader        EngagementType = "quiet_reader"
	EngagementSelectiveEngager   EngagementType = "selective_engager"
)

// Persona is the core soul object.
type Persona struct {
	ID              string         `json:"id"`
	Version         string         `json:"version"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DisplayName     string         `json:"display_name,omitempty"`
	Demographics    Demographics   `json:"demographics"`
	InterestGraph   InterestGraph  `json:"interest_graph,omitempty"`
	Language        Language       `json:"language,omitempty"`
	Circadian       Circadian      `json:"circadian"`
	PhysicalContext []string       `json:"physical_context_preferences,omitempty"`
	Engagement      Engagement     `json:"engagement"`
	Memory          Memory         `json:"memory,omitempty"`
	Visual          Visual         `json:"visual,omitempty"`
	DeviceProfileID string         `json:"device_profile_id,omitempty"`
	Status          Status         `json:"status"`
	CustomerOwned   bool           `json:"customer_owned"`
	TemplateID      *string        `json:"template_id,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
}

type Demographics struct {
	AgeRange           [2]int `json:"age_range"`
	GenderPresentation string `json:"gender_presentation,omitempty"`
	Location           string `json:"location"`
	LifeContext        string `json:"life_context,omitempty"`
}

type InterestGraph struct {
	Primary        []string `json:"primary,omitempty"`
	Secondary      []string `json:"secondary,omitempty"`
	AestheticStyle string   `json:"aesthetic_style,omitempty"`
}

type Language struct {
	PrimaryLanguage     string `json:"primary_language,omitempty"`
	DialectRegion       string `json:"dialect_region,omitempty"`
	EmojiHabits         string `json:"emoji_habits,omitempty"`
	VocabularyRegister  string `json:"vocabulary_register,omitempty"`
	CommentPhilosophy   string `json:"comment_philosophy,omitempty"`
}

type Circadian struct {
	TypicalWake  string   `json:"typical_wake,omitempty"`
	PeakActivity []string `json:"peak_activity,omitempty"`
	Timezone     string   `json:"timezone"`
}

type Engagement struct {
	Type             EngagementType `json:"type"`
	RiskTolerance    string         `json:"risk_tolerance,omitempty"`
	VolumePreference string         `json:"volume_preference,omitempty"`
}

type Memory struct {
	ShortTerm         []any     `json:"short_term,omitempty"`
	LongTermSummaries []any     `json:"long_term_summaries,omitempty"`
	LastUpdated       time.Time `json:"last_updated,omitempty"`
}

type Visual struct {
	ProfileAesthetic string `json:"profile_aesthetic,omitempty"`
	PhotoStyleHints  string `json:"photo_style_hints,omitempty"`
}

// New creates a minimal valid Persona.
func New(displayName, location, timezone string, ageMin, ageMax int, engagement EngagementType) *Persona {
	now := time.Now().UTC()
	return &Persona{
		ID:          uuid.New().String(),
		Version:     "0.1.0",
		CreatedAt:   now,
		UpdatedAt:   now,
		DisplayName: displayName,
		Demographics: Demographics{
			AgeRange: [2]int{ageMin, ageMax},
			Location: location,
		},
		Circadian: Circadian{
			Timezone: timezone,
		},
		Engagement: Engagement{
			Type: engagement,
		},
		Status:        StatusDraft,
		CustomerOwned: true,
		Metadata:      map[string]any{},
	}
}
