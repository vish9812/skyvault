-- Remove storage quota fields from profile table
drop index if exists profile_idx_storage_usage;
alter table profile drop column if exists storage_used_bytes;
alter table profile drop column if exists storage_quota_bytes;
