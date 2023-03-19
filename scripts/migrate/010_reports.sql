-- +migrate Up

-- 2023-03-19: Add moderation-related tables

alter table users add column is_admin boolean not null default false;

create table reports (
    id          serial primary key,
    -- we keep deleted users for 180 days after deletion, so it's fine to tie this to a user object
    user_id     text not null references users (id) on delete cascade,
    member_id   text null references members (id) on delete set null,
    reason      text not null,
    reporter_id text not null,

    created_at    timestamptz not null default now(),
    resolved_at   timestamptz,
    admin_id      text null references users (id) on delete set null,
    admin_comment text
);
