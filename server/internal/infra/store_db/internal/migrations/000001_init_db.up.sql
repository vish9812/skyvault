create table if not exists profiles (
	id bigserial primary key,
	email text not null unique,
	first_name text not null,
	last_name text not null,
	created_at timestamp not null default (timezone('utc', now())),
	updated_at timestamp not null default (timezone('utc', now()))
);

create table if not exists auths (
	id bigserial primary key,
	profile_id bigserial not null references profiles(id) on delete cascade,
	provider text not null,
	provider_user_id text not null,
	password_hash text null,
	created_at timestamp not null default (timezone('utc', now())),
	updated_at timestamp not null default (timezone('utc', now()))
);

create unique index if not exists auths_idx_unq_provider_user on auths(provider, provider_user_id);