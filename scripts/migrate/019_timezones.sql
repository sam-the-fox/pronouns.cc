-- +migrate Up

-- 2023-07-30: Add user timezones

alter table users add column timezone text null;
