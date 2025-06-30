CREATE TABLE IF NOT EXISTS users
(
    id UUID default gen_random_uuid() primary key,
    email varchar(20) not null,
    password varchar(100) not null,
    created_at timestamp not null default NOW()
);

CREATE TABLE IF NOT EXISTS sessions
(
    id serial primary key,
    user_id UUID not null references users(id) on delete cascade,
    refresh_token_hash text not null,
    user_agent text,
    ip_adress text,
    created_at timestamp not null default NOW(),
    expires_at timestamp not null,
    unique(refresh_token_hash)
);

CREATE INDEX IF NOT EXISTS idx_refresh_token_user_id ON sessions(user_id);