create table public.users (
    id bigserial primary key,
    created_at timestamp(0) with time zone not null default NOW(),
    name text not null,
    email citext unique not null,
    password_hash bytea not null
)
