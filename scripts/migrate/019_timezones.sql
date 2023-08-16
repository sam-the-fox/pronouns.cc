-- 2023-07-30: Add user timezones

-- +migrate Up

alter table users add column timezone text null;

-- +migrate Down

alter table users drop column timezone;
