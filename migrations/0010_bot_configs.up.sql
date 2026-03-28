CREATE TABLE IF NOT EXISTS bot_configs (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  token TEXT NOT NULL,
  database_id TEXT NOT NULL
    REFERENCES database_profiles(id)
    ON UPDATE RESTRICT
    ON DELETE RESTRICT,
  start_scenario TEXT NOT NULL DEFAULT 'reply_welcome',
  is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bot_configs_database_id ON bot_configs(database_id);
CREATE INDEX IF NOT EXISTS idx_bot_configs_is_enabled ON bot_configs(is_enabled);
