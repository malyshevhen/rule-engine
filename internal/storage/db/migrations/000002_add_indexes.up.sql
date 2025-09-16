-- Add indexes for performance optimization

-- Rules table indexes
CREATE INDEX idx_rules_enabled ON rules (enabled);
CREATE INDEX idx_rules_created_at ON rules (created_at DESC);

-- Triggers table indexes
CREATE INDEX idx_triggers_rule_id ON triggers (rule_id);
CREATE INDEX idx_triggers_type ON triggers (type);
CREATE INDEX idx_triggers_enabled ON triggers (enabled);

-- Actions table indexes (if needed for future queries)
CREATE INDEX idx_actions_enabled ON actions (enabled);

-- Execution logs indexes for monitoring and analytics
CREATE INDEX idx_execution_logs_rule_id ON execution_logs (rule_id);
CREATE INDEX idx_execution_logs_triggered_at ON execution_logs (triggered_at DESC);
CREATE INDEX idx_execution_logs_status ON execution_logs (execution_status);

-- Junction table indexes (already have primary keys, but these help with lookups)
CREATE INDEX idx_rule_triggers_trigger_id ON rule_triggers (trigger_id);
CREATE INDEX idx_rule_actions_action_id ON rule_actions (action_id);