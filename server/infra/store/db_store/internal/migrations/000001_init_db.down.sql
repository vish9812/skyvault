drop index if exists shares_unique_key;
drop table if exists shares;

drop index if exists files_deleted_at_idx;
drop index if exists files_unique_file_name_per_user;
drop table if exists files;

drop index if exists folders_deleted_at_idx;
drop index if exists folders_unique_folder_name_per_user;
drop table if exists folders;

drop table if exists users;
