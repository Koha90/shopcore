create table if not exists catalog_district_variants (
  id int generated always as identity primary key,
  district_id int not null references catalog_districts(id) on delete cascade,
  variant_id int not null references catalog_variants(id) on delete cascade,
  price integer not null,
  is_active boolean not null default true,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (district_id, variant_id)
);

create index if not exists idx_catalog_district_variants_district_id
    on catalog_district_variants(district_id);

create index if not exists idx_catalog_district_variants_variant_id
    on catalog_district_variants(variant_id);
