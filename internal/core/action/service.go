package action

import (
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
)

// Service handles business logic for actions
type Service struct {
	repo *actionStorage.Repository
}

// NewService creates a new action service
func NewService(repo *actionStorage.Repository) *Service {
	return &Service{repo: repo}
}

// TODO: Add methods for action execution
