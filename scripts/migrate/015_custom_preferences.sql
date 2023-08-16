-- 2023-04-19: Add custom preferences

-- +migrate Up

alter table users add column custom_preferences jsonb not null default '{}';

-- +migrate Down

alter table users drop column custom_preferences;
