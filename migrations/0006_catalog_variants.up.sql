create table if not exists catalog_variants (
  id int generated always as identity primary key,
  product_id int not null references catalog_products(id) on delete cascade,
  code text not null,
  name text not null,
  name_latin text not null default '',
  description text not null default '',
  image_url text not null default '',
  is_active boolean not null default true,
  sort_order integer not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (product_id, code),
  unique (product_id, name)
);

create index if not exists idx_catalog_variants_product_id
    on catalog_variants(product_id);
