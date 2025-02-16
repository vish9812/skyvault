create table if not exists profile (
	id bigserial primary key,
	email text not null unique,
	full_name text not null,
	avatar bytea,
	created_at timestamp not null default (timezone('utc', now())),
	updated_at timestamp not null default (timezone('utc', now()))
);

create table if not exists auth (
	id bigserial primary key,
	profile_id bigint not null references profile(id) on delete cascade,
	provider text not null,
	provider_user_id text not null,
	password_hash text null, -- null if provider is not 'email'
	created_at timestamp not null default (timezone('utc', now())),
	updated_at timestamp not null default (timezone('utc', now()))
);

create index if not exists auth_idx_profile_id on auth(profile_id);
create unique index if not exists auth_idx_unq_provider_user on auth(provider, provider_user_id);

create table if not exists folder_info (
	id bigserial primary key,
	owner_id bigint not null references profile(id) on delete cascade,
	name text not null,
	parent_folder_id bigint references folder_info(id) on delete cascade,
	trashed_at timestamp,
	created_at timestamp not null default (timezone('utc', now())),
	updated_at timestamp not null default (timezone('utc', now()))
);

create unique index if not exists folder_info_idx_unq_folder_per_user 
on folder_info(owner_id, parent_folder_id, name) nulls not distinct
where trashed_at is null;

create table if not exists file_info (
	id bigserial primary key,
	owner_id bigint not null references profile(id) on delete cascade,
	folder_id bigint references folder_info(id) on delete cascade,
	name text not null,
	generated_name text not null,
	size bigint not null,
	extension text,
	mime_type text not null,
	category text,
	preview bytea,
	trashed_at timestamp,
	created_at timestamp not null default (timezone('utc', now())),
	updated_at timestamp not null default (timezone('utc', now()))
);

create unique index if not exists file_info_idx_unq_file_per_user 
on file_info(owner_id, folder_id, name) nulls not distinct
where trashed_at is null;

create index if not exists file_info_idx_category
on file_info(owner_id, category)
where trashed_at is null;

-- Contacts list for users
create table if not exists contact (
    id bigserial primary key,
    owner_id bigint not null references profile(id) on delete cascade,
    email text not null,
    name text,
    created_at timestamp not null default (timezone('utc', now())),
    updated_at timestamp not null default (timezone('utc', now()))
);

create unique index if not exists contact_idx_unq_email_per_user 
on contact(owner_id, email);

-- Sharing configuration for resources (files/folders)
create table if not exists share_config (
    id bigserial primary key,
    owner_id bigint not null references profile(id) on delete cascade,
    resource_type text not null check (resource_type in ('file', 'folder')),
    resource_id bigint not null,
    password_hash text, -- optional password protection
    max_downloads int, -- null means unlimited
    expires_at timestamp, -- null means never expires
    created_at timestamp not null default (timezone('utc', now())),
    updated_at timestamp not null default (timezone('utc', now()))
);

-- Ensure resource_id points to correct table based on resource_type
create index if not exists share_config_idx_file 
on share_config(resource_id)
where resource_type = 'file';

create index if not exists share_config_idx_folder
on share_config(resource_id) 
where resource_type = 'folder';

-- Share recipients
create table if not exists share_recipient (
    id bigserial primary key,
    share_config_id bigint not null references share_config(id) on delete cascade,
    email text not null,
    downloads_count int not null default 0,
    created_at timestamp not null default (timezone('utc', now())),
    updated_at timestamp not null default (timezone('utc', now()))
);

-- For looking up shares by email
create index if not exists share_recipient_idx_email
on share_recipient(email);

-- Ensure unique recipient per share
create unique index if not exists share_recipient_idx_unq_email
on share_recipient(share_config_id, email);

-- Track share access history
create table if not exists share_access (
    id bigserial primary key,
    share_id bigint not null references share(id) on delete cascade,
    accessed_at timestamp not null default (timezone('utc', now())),
    accessed_from_ip text not null
);
