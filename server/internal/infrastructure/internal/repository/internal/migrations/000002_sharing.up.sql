create table if not exists contact (
    id bigserial primary key,
    owner_id bigint not null references profile(id) on delete cascade,
    email text not null,
    name text not null,
    created_at timestamp not null default (timezone('utc', now())),
    updated_at timestamp not null default (timezone('utc', now()))
);

create unique index if not exists contact_idx_unq_email_per_user 
on contact(owner_id, email);

create unique index if not exists contact_idx_unq_email_per_user_name
on contact(owner_id, email, name);

create table if not exists contact_group (
    id bigserial primary key,
    owner_id bigint not null references profile(id) on delete cascade,
    name text not null,
    created_at timestamp not null default (timezone('utc', now())),
    updated_at timestamp not null default (timezone('utc', now()))
);

create unique index if not exists contact_group_idx_unq_name
on contact_group(owner_id, name);

create table if not exists contact_group_member (
    id bigserial primary key,
    group_id bigint not null references contact_group(id) on delete cascade,
    contact_id bigint not null references contact(id) on delete cascade,
    created_at timestamp not null default (timezone('utc', now()))
);

create unique index if not exists contact_group_member_idx_unq_contact
on contact_group_member(group_id, contact_id);

create table if not exists share_config (
    id bigserial primary key,
    custom_id uuid not null, -- for url path
    owner_id bigint not null references profile(id) on delete cascade,
    file_id bigint references file_info(id) on delete cascade,
    folder_id bigint references folder_info(id) on delete cascade,
    password_hash text, -- optional password protection
    max_downloads bigint, -- null means unlimited
    expires_at timestamp, -- null means never expires
    created_at timestamp not null default (timezone('utc', now())),
    updated_at timestamp not null default (timezone('utc', now())),
    constraint either_file_or_folder
        check (
            (file_id is not null and folder_id is null) or
            (file_id is null and folder_id is not null)
        )
);

create unique index if not exists share_config_idx_unq_custom_id
on share_config(custom_id);

create unique index if not exists share_config_idx_unq_file 
on share_config(file_id);

create unique index if not exists share_config_idx_unq_folder
on share_config(folder_id);

create table if not exists share_recipient (
    id bigserial primary key,
    share_config_id bigint not null references share_config(id) on delete cascade,
    contact_id bigint references contact(id) on delete cascade,
    group_id bigint references contact_group(id) on delete cascade,
    email text,
    downloads_count bigint not null default 0,
    created_at timestamp not null default (timezone('utc', now())),
    updated_at timestamp not null default (timezone('utc', now())),
    constraint either_contact_or_group_or_email
        check (
            (contact_id is not null and group_id is null and email is null) or
            (group_id is not null and contact_id is null and email is null) or
            (email is not null and contact_id is null and group_id is null)
        )
);

create unique index if not exists share_recipient_idx_unq_config
on share_recipient(share_config_id);

create unique index if not exists share_recipient_idx_unq_contact
on share_recipient(share_config_id, contact_id);

create unique index if not exists share_recipient_idx_unq_group
on share_recipient(share_config_id, group_id);

create unique index if not exists share_recipient_idx_unq_email
on share_recipient(share_config_id, email);
