drop index if exists auth_idx_unq_provider_user;
drop table if exists auth;

drop index if exists file_info_idx_unq_file_per_user;
drop index if exists file_info_idx_category;
drop table if exists file_info;

drop index if exists folder_info_idx_unq_folder_per_user;
drop table if exists folder_info;

drop table if exists profile;
