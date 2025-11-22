-- Add upload_session table to track chunked upload sessions and prevent quota bypass attacks
create table if not exists upload_session (
    id uuid primary key default gen_random_uuid(),
    owner_id uuid not null references profile(id) on delete cascade,
    file_name text not null,
    file_size bigint not null check (file_size > 0),
    mime_type text not null,
    folder_id uuid references folder_info(id) on delete set null,
    total_chunks int not null check (total_chunks > 0),
    quota_allocated bigint not null check (quota_allocated > 0),
    created_at timestamp not null default now(),
    expires_at timestamp not null
);

-- Index for querying by owner
create index if not exists idx_upload_session_owner_id on upload_session(owner_id);

-- Index for cleanup queries (finding expired sessions)
create index if not exists idx_upload_session_expires_at on upload_session(expires_at);
