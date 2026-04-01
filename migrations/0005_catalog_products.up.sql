create table if not exists catalog_products (
  id int generated always as identity primary key,
  category_id int not null references catalog_categories(id) on delete restrict,
  district_id int not null references catalog_districts(id) on delete cascade,
  code text not null,
  name text not null,
  name_latin text not null default '',
  description text not null default '',
  is_active boolean not null default true,
  sort_order integer not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (district_id, category_id, code)
);

create index if not exists idx_catalog_products_category_id
    on catalog_products(category_id);

create index if not exists idx_catalog_products_category_id
    on catalog_products(district_id);
