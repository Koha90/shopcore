create table if not exists catalog_districts (
  id int generated always as identity primary key,
  city_id int not null references cities(id) on delete cascade,
  code text not null,
  name text not null,
  name_latin text not null default '',
  is_active boolean not null default true,
  sort_order integer not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (city_id, code),
  unique (city_id, name)
);

create index if not exists idx_catalog_districts_city_id
    on catalog_districts(city_id);
