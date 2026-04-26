create table if not exists bot_configs (
  id text primary key,
  name text not null,
  token text not null,

  database_id text not null
    references database_profiles(id)
    on update restrict
    on delete restrict,

  start_scenario text not null default 'reply_welcome',

  telegram_admin_user_ids bigint[] not null default '{}',
  admin_orders_chat_id bigint not null default 0,

  telegram_bot_id bigint not null default 0,
  telegram_username text not null default '',
  telegram_bot_name text not null default '',

  is_enabled boolean not null default true,
  updated_at timestamptz not null default now()
);

create index if not exists idx_bot_configs_database_id on bot_configs(database_id);
create index if not exists idx_bot_configs_is_enabled on bot_configs(is_enabled);
