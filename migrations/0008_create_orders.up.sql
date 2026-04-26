create table if not exists orders (
  id bigint generated always as identity primary key,

  bot_id text not null,
  bot_name text not null,

  chat_id bigint not null,
  user_id bigint not null,
  user_name text not null default '',
  user_username text not null default '',

  city_id text not null,
  city_name text not null,

  district_id text not null,
  district_name text not null,

  product_id text not null,
  product_name text not null,

  variant_id text not null,
  variant_name text not null,

  price_text text not null default '',

  status text not null,
  created_at timestamptz not null default now()
);

create index if not exists idx_orders_created_at_desc
    on orders (created_at desc);

create index if not exists idx_orders_status_created_at_desc
    on orders (status, created_at desc);

create index if not exists idx_orders_bot_id_created_at_desc
    on orders (bot_id, created_at desc);
