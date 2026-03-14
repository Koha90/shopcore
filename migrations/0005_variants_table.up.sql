CREATE TABLE IF NOT EXISTS product_variants (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  district_id BIGINT NOT NULL REFERENCES districts(id) ON DELETE RESTRICT,
  pack_size TEXT NOT NULL,
  price BIGINT NOT NULL CHECK (price > 0),
  archived_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_variants_product_id ON product_variants(product_id);
CREATE INDEX IF NOT EXISTS idx_variants_district_id ON product_variants(district_id);

-- Один активный вариант на product + district + pack_size.
CREATE UNIQUE INDEX IF NOT EXISTS uq_variants_active_product_district_pack
  ON product_variants(product_id, district_id, pack_size)
  WHERE archived_at IS NULL;
