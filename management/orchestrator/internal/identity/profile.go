package identity

import (
	"time"

	"github.com/google/uuid"
)

// DeviceProfile represents a coherent Android device identity
// that will be bound to a Persona.
type DeviceProfile struct {
	ID            string            `json:"id"`
	CreatedAt     time.Time         `json:"created_at"`
	Model         string            `json:"model"`
	Manufacturer  string            `json:"manufacturer"`
	AndroidVersion string          `json:"android_version"`
	BuildFingerprint string        `json:"build_fingerprint,omitempty"`
	AndroidID     string            `json:"android_id"`
	GAID          string            `json:"gaid,omitempty"` // Google Advertising ID
	Serial        string            `json:"serial,omitempty"`
	Properties    map[string]string `json:"properties,omitempty"` // build.prop style
	Notes         string            `json:"notes,omitempty"`
}

// NewBasicProfile creates a simple coherent device profile.
// Later this will become much richer and persona-aware.
func NewBasicProfile(model, manufacturer, androidVersion string) *DeviceProfile {
	return &DeviceProfile{
		ID:             uuid.New().String(),
		CreatedAt:      time.Now().UTC(),
		Model:          model,
		Manufacturer:   manufacturer,
		AndroidVersion: androidVersion,
		AndroidID:      uuid.New().String(),
		GAID:           uuid.New().String(),
		Serial:         uuid.New().String()[:12],
		Properties:     map[string]string{},
	}
}

// DefaultProfiles returns a small set of reasonable starting profiles.
func DefaultProfiles() []*DeviceProfile {
	return []*DeviceProfile{
		NewBasicProfile("Pixel 7", "Google", "14"),
		NewBasicProfile("Pixel 6a", "Google", "14"),
		NewBasicProfile("Galaxy S23", "Samsung", "14"),
		NewBasicProfile("OnePlus 11", "OnePlus", "14"),
	}
}
