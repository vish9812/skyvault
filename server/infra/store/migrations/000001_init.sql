create table if not exists users (
	id uuid primary key,
	first_name text not null,
	last_name text not null,
	username text not null unique,
	email text not null unique,
	password_hash text not null,
	auth_service text not null,
	created_at timestamp not null default (timezone('utc', now())),
	updated_at timestamp not null default (timezone('utc', now()))
);

create table if not exists folders (
	id uuid primary key,
	user_id uuid not null references users(id) on delete cascade,
	name text not null,
	parent_folder_id uuid references folders(id) on delete cascade,
	created_at timestamp not null default (timezone('utc', now())),
	updated_at timestamp not null default (timezone('utc', now())),
	deleted_at timestamp
);

create unique index if not exists folders_unique_folder_name_per_user 
on folders(user_id, parent_folder_id, name)
where deleted_at is null;

create index if not exists folders_deleted_at_idx on folders(deleted_at) where deleted_at is null;

create table if not exists files (
	id uuid primary key,
	user_id uuid not null references users(id) on delete cascade,
	folder_id uuid references folders(id) on delete cascade,
	name text not null,
	file_size_bytes bigint not null,
	created_at timestamp not null default (timezone('utc', now())),
	updated_at timestamp not null default (timezone('utc', now())),
	deleted_at timestamp
);

create unique index if not exists files_unique_file_name_per_user 
on files(user_id, folder_id, name)
where deleted_at is null;

create index if not exists files_deleted_at_idx on files(deleted_at) where deleted_at is null;

create table if not exists shares (
	id uuid primary key,
	owner_id uuid not null references users(id) on delete cascade,
	recipient_id uuid references users(id) on delete cascade,
	folder_id uuid references folders(id) on delete cascade,
	file_id uuid references files(id) on delete cascade,
	permission text not null,
	shared_at timestamp not null default (timezone('utc', now())),
	shared_until timestamp
);

create unique index if not exists shares_unique_key on shares(owner_id, recipient_id, folder_id, file_id);
