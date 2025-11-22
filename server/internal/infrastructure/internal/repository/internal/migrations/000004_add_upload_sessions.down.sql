-- Drop upload_session table and related indices
drop index if exists idx_upload_session_expires_at;
drop index if exists idx_upload_session_owner_id;
drop table if exists upload_session;
