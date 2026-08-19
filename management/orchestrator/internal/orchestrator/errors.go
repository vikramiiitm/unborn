package orchestrator

import "errors"

var (
	ErrMaxInstancesReached = errors.New("maximum number of instances reached")
	ErrInstanceNotFound    = errors.New("instance not found")
	ErrPersonaNotFound     = errors.New("persona not found")
)
