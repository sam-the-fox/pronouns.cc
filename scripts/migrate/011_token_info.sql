-- +migrate Up

-- 2023-03-30: Add token information to database

alter table tokens add column api_only boolean not null default false;
alter table tokens add column read_only boolean not null default false;
