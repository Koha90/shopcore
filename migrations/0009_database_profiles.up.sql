CREATE TABLE IF NOT EXISTS database_profiles (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  driver TEXT NOT NULL,
  dsn TEXT NOT NULL,
  is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_database_profiles_is_enabled ON database_profiles(is_enabled);
