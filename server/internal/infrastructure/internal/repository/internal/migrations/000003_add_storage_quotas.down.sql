-- Remove storage quota fields from profile table
alter table profile drop column if exists storage_used;
alter table profile drop column if exists storage_quota;
