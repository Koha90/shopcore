create table if not exists catalog_variants (
  id int generated always as identity primary key,
  product_id int not null references catalog_products(id) on delete cascade,
  code text not null,
  name text not null,
  name_latin text not null default '',
  discription text not null default '',
  price_minor bigint not null,
  currency_code char(3) not null default 'RUB',
  is_active boolean not null default true,
  sort_order integer not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (product_id, code)
);

create index if not exists idx_catalog_variants_product_id
    on catalog_variants(product_id);
