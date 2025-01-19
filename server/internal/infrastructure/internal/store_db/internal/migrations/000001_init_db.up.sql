create table if not exists profiles (
	id bigserial primary key,
	email text not null unique,
	full_name text not null,
	avatar bytea,
	created_at timestamp not null default (timezone('utc', now())),
	updated_at timestamp not null default (timezone('utc', now()))
);

create table if not exists auths (
	id bigserial primary key,
	profile_id bigint not null references profiles(id) on delete cascade,
	provider text not null,
	provider_user_id text not null,
	password_hash text null, -- null if provider is not 'email'
	created_at timestamp not null default (timezone('utc', now())),
	updated_at timestamp not null default (timezone('utc', now()))
);

create unique index if not exists auths_idx_unq_provider_user on auths(provider, provider_user_id);

create table if not exists folders (
	id bigserial primary key,
	owner_id bigint not null references profiles(id) on delete cascade,
	name text not null,
	parent_folder_id bigint references folders(id) on delete cascade,
	trashed_at timestamp,
	created_at timestamp not null default (timezone('utc', now())),
	updated_at timestamp not null default (timezone('utc', now()))
);

create unique index if not exists folders_idx_unq_folder_per_user 
on folders(owner_id, parent_folder_id, name)
where trashed_at is null;

create table if not exists files (
	id bigserial primary key,
	owner_id bigint not null references profiles(id) on delete cascade,
	folder_id bigint references folders(id) on delete cascade,
	name text not null,
	size_bytes bigint not null,
	extension text,
	mime_type text not null,
	trashed_at timestamp,
	created_at timestamp not null default (timezone('utc', now())),
	updated_at timestamp not null default (timezone('utc', now()))
);

create unique index if not exists files_idx_unq_file_per_user 
on files(owner_id, folder_id, name)
where trashed_at is null;
