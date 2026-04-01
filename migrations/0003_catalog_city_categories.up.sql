create table if not exists catalog_city_categories (
  city_id int not null references cities(id) on delete cascade,
  category_id int not null references catalog_categories(id) on delete cascade,
  sort_order integer not null default 0,
  created_at timestamptz not null default now(),
  primary key (city_id, category_id)
);

create index if not exists idx_catalog_city_categories_category_id
    on catalog_city_categories(category_id);
