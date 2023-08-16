-- 2023-03-30: Add token information to database

-- +migrate Up

alter table tokens add column api_only boolean not null default false;
alter table tokens add column read_only boolean not null default false;

-- +migrate Down

alter table tokens drop column api_only;
alter table tokens drop column read_only;
