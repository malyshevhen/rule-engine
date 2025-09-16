package context

// Service manages execution contexts
type Service struct {
	// TODO: Add dependencies like cache, etc.
}

// NewService creates a new execution context service
func NewService() *Service {
	return &Service{}
}

// CreateContext creates a new execution context for a rule execution
func (s *Service) CreateContext(ruleID, triggerID string) *ExecutionContext {
	return &ExecutionContext{
		RuleID:    ruleID,
		TriggerID: triggerID,
		Data:      make(map[string]any),
	}
}
