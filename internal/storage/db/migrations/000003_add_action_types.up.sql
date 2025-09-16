ALTER TABLE actions ADD COLUMN type VARCHAR(50) NOT NULL DEFAULT 'lua_script';
ALTER TABLE actions ADD COLUMN params TEXT;

-- Migrate existing lua_script to params as JSON
UPDATE actions
SET params = JSON_BUILD_OBJECT('script', lua_script)::TEXT
WHERE params IS NULL;

-- Make params NOT NULL after migration
ALTER TABLE actions ALTER COLUMN params SET NOT NULL;
