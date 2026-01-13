-- Add storage quota fields to profile table
alter table profile add column storage_quota bigint not null default 0 check (storage_quota >= 0);
alter table profile add column storage_used bigint not null default 0 check (storage_used >= 0) ;

-- Check constraint for storage_used not exceeding storage_quota
alter table profile add constraint storage_used_within_quota check (storage_used <= storage_quota

