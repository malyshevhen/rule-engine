package trigger

import (
	"github.com/malyshevhen/rule-engine/internal/storage/trigger"
)

// Service handles business logic for triggers
type Service struct {
	repo *trigger.Repository
}

// NewService creates a new trigger service
func NewService(repo *trigger.Repository) *Service {
	return &Service{repo: repo}
}

// TODO: Add methods for trigger evaluation
