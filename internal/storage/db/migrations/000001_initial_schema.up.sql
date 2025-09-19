-- Create custom types
CREATE TYPE trigger_type AS ENUM ('CONDITIONAL', 'CRON');
CREATE TYPE execution_status AS ENUM ('SUCCESS', 'FAILURE', 'TIMEOUT');

-- Rules table
CREATE TABLE rules (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    lua_script TEXT NOT NULL,
    priority   INTEGER NOT NULL DEFAULT 0,
    enabled    BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Actions table
CREATE TABLE actions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type       VARCHAR(50) NOT NULL DEFAULT 'lua_script',
    params     TEXT NOT NULL, -- JSON string for parameters
    enabled    BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Triggers table
CREATE TABLE triggers (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id          UUID NOT NULL REFERENCES rules (id) ON DELETE CASCADE,
    type             TRIGGER_TYPE NOT NULL,
    condition_script TEXT NOT NULL,
    enabled          BOOLEAN NOT NULL DEFAULT true,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Junction table for rule-actions many-to-many
CREATE TABLE rule_actions (
    rule_id   UUID NOT NULL REFERENCES rules (id) ON DELETE CASCADE,
    action_id UUID NOT NULL REFERENCES actions (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (rule_id, action_id)
);

-- Execution logs table
CREATE TABLE execution_logs (
    id               BIGSERIAL PRIMARY KEY,
    rule_id          UUID NOT NULL REFERENCES rules (id) ON DELETE CASCADE,
    triggered_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    execution_status EXECUTION_STATUS NOT NULL,
    duration_ms      INTEGER,
    output_log       TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes for performance
CREATE INDEX idx_rules_enabled ON rules (enabled);
CREATE INDEX idx_rules_priority ON rules (priority);
CREATE INDEX idx_actions_enabled ON actions (enabled);
CREATE INDEX idx_triggers_rule_id ON triggers (rule_id);
CREATE INDEX idx_triggers_enabled ON triggers (enabled);
CREATE INDEX idx_execution_logs_rule_id ON execution_logs (rule_id);
CREATE INDEX idx_execution_logs_triggered_at ON execution_logs (triggered_at);
CREATE INDEX idx_rule_actions_rule_id ON rule_actions (rule_id);
