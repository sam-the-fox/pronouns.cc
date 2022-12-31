-- +migrate Up

-- 2022-12-23: Add database-backed tokens
create table tokens (
    user_id     text not null references users (id) on delete cascade,
    token_id    text primary key,
    invalidated boolean not null default false,
    created     timestamptz not null default now(),
    expires     timestamptz not null
);

-- Unrelatedly, this migration also changes the column type for invites.created to timestamptz (from plain timestamp)
-- This does not change anything code-wise, but it's recommended over plain timestamp because plain timestamp does not handle timezones correctly
alter table invites alter column created type timestamptz;
alter table invites alter column created set default now();
