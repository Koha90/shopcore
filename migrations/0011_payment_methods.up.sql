create table if not exists payment_methods (
  id int generated always as identity primary key,
  code text not null unique,
  name text not null,
  kind text not null,
  is_active boolean not null default true,
  sort_order integer not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  constraint payment_methods_code_not_empty check (btrim(code) <> ''),
  constraint payment_methods_name_not_empty check (btrim(name) <> ''),
  constraint payment_methods_kind_allowed check (
    kind in (
      'bank_card',
      'sbp',
      'mobile_phone',
      'btc',
      'eth',
      'cash',
      'manual'
    )
  )
);

create index if not exists idx_payment_methods_active_sort
  on payment_methods(is_active, sort_order, id);

create table if not exists bot_payment_methods (
  id int generated always as identity primary key,
  bot_id text not null,
  payment_methods_id int not null references payment_methods(id) on delete restrict,
  display_name text not null default '',
  is_active boolean not null default true,

  -- Basic points:
  -- 100 = 1%
  -- 250 = 2.5%
  -- 10000 = 100%
  extra_percent_bps integer not null default 0,

  sort_order integer not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  unique (bot_id, payment_methods_id),

  constraint bot_payment_methods_bot_id_empty check (btrim(bot_id) <> ''),
  constraint bot_payment_methods_extra_percent_bps_non_negative check (extra_percent_bps >= 0)
);

create index if not exists idx_bot_payment_methods_bot_active_sort
  on bot_payment_methods(bot_id, is_active, sort_order, id)
