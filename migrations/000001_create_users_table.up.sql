create table public.users (
    id bigserial primary key,
    created_at timestamp(0) with time zone not null default NOW(),
    name text unique not null,
    email citext not null,
    password_hash bytea not null
)
