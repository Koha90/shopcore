create table if not exists catalog_categories (
  id int generated always as identity primary key,
  code text not null unique,
  name text not null unique,
  name_latin text not null default '',
  description text not null default '',
  is_active boolean not null default true,
  sort_order integer not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
);
