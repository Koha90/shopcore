CREATE TABLE IF NOT EXISTS users (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

  tg_id BIGINT UNIQUE,
  tg_name TEXT,

  email TEXT UNIQUE,
  password_hash TEXT,

  role TEXT NOT NULL CHECK (role IN ('customer', 'admin')),
  balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),

  is_enabled BOOLEAN NOT NULL DEFAULT FALSE,

  admin_access_expires_at TIMESTAMPTZ NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CHECK (
    tg_id IS NOT NULL
    OR (email IS NOT NULL AND password_hash IS NOT NULL)
  )
);

CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_enable ON users(is_enable);
CREATE INDEX IF NOT EXISTS idx_users_tg_id ON users(tg_id);
