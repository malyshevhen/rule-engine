CREATE TABLE rules (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    lua_script TEXT NOT NULL,
    enabled    BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TYPE trigger_type AS ENUM ('CONDITIONAL', 'CRON');

CREATE TABLE triggers (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id          UUID NOT NULL REFERENCES rules (id) ON DELETE CASCADE,
    type             TRIGGER_TYPE NOT NULL,
    condition_script TEXT NOT NULL,
    enabled          BOOLEAN NOT NULL DEFAULT true,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE actions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lua_script TEXT NOT NULL,
    enabled    BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TYPE execution_status AS ENUM ('SUCCESS', 'FAILURE', 'TIMEOUT');

CREATE TABLE execution_logs (
    id               BIGSERIAL PRIMARY KEY,
    rule_id          UUID NOT NULL REFERENCES rules (id),
    triggered_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    execution_status VARCHAR(20) NOT NULL,
    duration_ms      INT,
    output_log       TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE rule_triggers (
    rule_id    UUID NOT NULL REFERENCES rules (id) ON DELETE CASCADE,
    trigger_id UUID NOT NULL REFERENCES triggers (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (rule_id, trigger_id)
);

ALTER TABLE rule_triggers ADD CONSTRAINT rule_triggers_unique UNIQUE (
    rule_id, trigger_id
);

CREATE TABLE rule_actions (
    rule_id    UUID NOT NULL REFERENCES rules (id) ON DELETE CASCADE,
    action_id  UUID NOT NULL REFERENCES actions (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (rule_id, action_id)
);

ALTER TABLE rule_actions ADD CONSTRAINT rule_actions_unique UNIQUE (
    rule_id, action_id
);
