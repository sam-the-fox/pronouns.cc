-- 2023-04-18: Add Google oauth

-- +migrate Up

alter table users add column google text null;
alter table users add column google_username text null;

-- +migrate Down

alter table users drop column google;
alter table users drop column google_username;
