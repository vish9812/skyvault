-- Add storage quota fields to profile table
alter table profile add column storage_quota bigint not null default 0;
alter table profile add column storage_used bigint not null default 0;

-- Create index for efficient quota queries
create index if not exists profile_idx_storage_usage on profile(storage_used);
