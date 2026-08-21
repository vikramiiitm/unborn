package license

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"time"
)

// Info is the decoded license payload (Phase 1: HMAC-signed JSON, offline).
type Info struct {
	CustomerID   string    `json:"customer_id"`
	MaxInstances int       `json:"max_instances"`
	Tier         string    `json:"tier"` // starter, growth, scale, dev
	ExpiresAt    time.Time `json:"expires_at"`
	IssuedAt     time.Time `json:"issued_at"`
}

// Status is what the API returns.
type Status struct {
	Valid        bool      `json:"valid"`
	CustomerID   string    `json:"customer_id,omitempty"`
	MaxInstances int       `json:"max_instances"`
	Tier         string    `json:"tier,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	Reason       string    `json:"reason,omitempty"`
	DevMode      bool      `json:"dev_mode"`
}

// Service holds the active license.
type Service struct {
	mu     sync.RWMutex
	info   *Info
	secret []byte
}

func NewService() *Service {
	secret := os.Getenv("UNBORN_LICENSE_SECRET")
	if secret == "" {
		secret = "unborn-dev-secret-change-me"
	}
	s := &Service{secret: []byte(secret)}
	// Dev default: generous limits until a key is activated
	s.info = &Info{
		CustomerID:   "dev",
		MaxInstances: 10,
		Tier:         "dev",
		IssuedAt:     time.Now().UTC(),
		ExpiresAt:    time.Now().UTC().Add(365 * 24 * time.Hour),
	}
	if key := os.Getenv("UNBORN_LICENSE_KEY"); key != "" {
		_ = s.Activate(key)
	}
	return s
}

func (s *Service) Status() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.info == nil {
		return Status{Valid: false, Reason: "no license", MaxInstances: 0}
	}
	if time.Now().UTC().After(s.info.ExpiresAt) {
		return Status{Valid: false, Reason: "expired", CustomerID: s.info.CustomerID, MaxInstances: 0, Tier: s.info.Tier, ExpiresAt: s.info.ExpiresAt}
	}
	return Status{
		Valid:        true,
		CustomerID:   s.info.CustomerID,
		MaxInstances: s.info.MaxInstances,
		Tier:         s.info.Tier,
		ExpiresAt:    s.info.ExpiresAt,
		DevMode:      s.info.Tier == "dev",
	}
}

func (s *Service) MaxInstances() int {
	st := s.Status()
	if !st.Valid {
		return 0
	}
	return st.MaxInstances
}

// Activate accepts payload.hex_signature
func (s *Service) Activate(key string) error {
	parts := strings.SplitN(key, ".", 2)
	if len(parts) != 2 {
		return errors.New("invalid license key format")
	}
	payload, sig := parts[0], parts[1]
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(sig)) {
		return errors.New("invalid license signature")
	}
	var info Info
	if err := json.Unmarshal([]byte(payload), &info); err != nil {
		// try hex-encoded payload
		raw, err2 := hex.DecodeString(payload)
		if err2 != nil {
			return errors.New("invalid license payload")
		}
		if err := json.Unmarshal(raw, &info); err != nil {
			return errors.New("invalid license payload json")
		}
	}
	if info.MaxInstances <= 0 {
		info.MaxInstances = 5
	}
	s.mu.Lock()
	s.info = &info
	s.mu.Unlock()
	return nil
}

// IssueDevKey builds a signed key for local testing (same secret).
func (s *Service) IssueDevKey(customerID string, max int, tier string, days int) string {
	if max <= 0 {
		max = 10
	}
	if days <= 0 {
		days = 30
	}
	if tier == "" {
		tier = "starter"
	}
	info := Info{
		CustomerID:   customerID,
		MaxInstances: max,
		Tier:         tier,
		IssuedAt:     time.Now().UTC(),
		ExpiresAt:    time.Now().UTC().Add(time.Duration(days) * 24 * time.Hour),
	}
	raw, _ := json.Marshal(info)
	payload := hex.EncodeToString(raw)
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))
	return payload + "." + sig
}
