CREATE TABLE IF NOT EXISTS districts (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  city_id INT NOT NULL REFERENCES cities(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  UNIQUE(city_id, name)
);

CREATE INDEX IF NOT EXISTS idx_districts_city_id ON districts(city_id);
