-- Drop tables in reverse order to handle foreign keys
DROP TABLE IF EXISTS execution_logs;
DROP TABLE IF EXISTS rule_actions;
DROP TABLE IF EXISTS triggers;
DROP TABLE IF EXISTS actions;
DROP TABLE IF EXISTS rules;

-- Drop custom types
DROP TYPE IF EXISTS execution_status;
DROP TYPE IF EXISTS trigger_type;