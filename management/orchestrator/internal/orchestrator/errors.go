package orchestrator

import "errors"

var (
	ErrMaxInstancesReached = errors.New("maximum number of instances reached")
	ErrInstanceNotFound    = errors.New("instance not found")
	ErrPersonaNotFound     = errors.New("persona not found")
	ErrLicenseInvalid      = errors.New("license invalid or expired")
	ErrSimulatedNoScreen   = errors.New("simulated body has no screen — start a real Redroid body")
)
