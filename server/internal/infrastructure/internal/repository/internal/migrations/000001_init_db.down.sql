drop table if exists share_access;

drop index if exists share_recipient_idx_unq_email;
drop index if exists share_recipient_idx_email;
drop table if exists share_recipient;

drop index if exists share_config_idx_folder;
drop index if exists share_config_idx_file;
drop table if exists share_config;

drop index if exists contact_group_member_idx_unq_contact;
drop table if exists contact_group_member;

drop index if exists contact_group_idx_unq_name;
drop table if exists contact_group;

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
