drop table if exists share_access;

drop index if exists share_idx_email;
drop index if exists share_idx_folder;
drop index if exists share_idx_file;
drop table if exists share;

drop index if exists contact_idx_unq_email_per_user;
drop table if exists contact;

drop index if exists auth_idx_profile_id;
drop index if exists auth_idx_unq_provider_user;
drop table if exists auth;

drop index if exists file_info_idx_unq_file_per_user;
drop index if exists file_info_idx_category;
drop table if exists file_info;

drop index if exists folder_info_idx_unq_folder_per_user;
drop table if exists folder_info;

drop table if exists profile;
