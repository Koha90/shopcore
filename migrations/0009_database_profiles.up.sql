create table if not exists database_profiles (
  id text primary key,
  name text not null unique,
  driver text not null,
  dsn text not null,
  is_enabled boolean not null default true,
  updated_at timestamptz not null default now()
);

create index if not exists idx_database_profiles_is_enabled on database_profiles(is_enabled);
