-- +migrate Up

-- 2023-03-11: Change composite type arrays to use jsonb columns
-- Composite types aren't actually supported by pgx and this allows us to drop pggen as a dev dependency.

-- Delete old columns
alter table users drop column names;
alter table users drop column pronouns;

alter table members drop column names;
alter table members drop column pronouns;

alter table user_fields drop column entries;
alter table member_fields drop column entries;

-- Create new columns
alter table users add column names jsonb not null default '[]';
alter table users add column pronouns jsonb not null default '[]';

alter table members add column names jsonb not null default '[]';
alter table members add column pronouns jsonb not null default '[]';

alter table user_fields add column entries jsonb not null default '[]';
alter table member_fields add column entries jsonb not null default '[]';
