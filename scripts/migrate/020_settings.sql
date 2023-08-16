-- 2023-08-16: Server-side user settings storage (for now, just whether specific notices are dismissed)

-- +migrate Up

alter table users add column settings jsonb not null default '{}';

-- +migrate Down

alter table users drop column settings;
