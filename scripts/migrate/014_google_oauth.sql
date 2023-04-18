-- +migrate Up

-- 2023-04-18: Add Google oauth

alter table users add column google text null;
alter table users add column google_username text null;
