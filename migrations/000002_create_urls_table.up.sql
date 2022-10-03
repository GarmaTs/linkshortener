create table public.urls (
	id bigserial primary key,
	user_id int,
    created_at timestamp(0) with time zone not null default NOW(),
	short_url text,
	full_url text,
	unique(user_id, full_url)
)
