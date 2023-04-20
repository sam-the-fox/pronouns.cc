-- +migrate Up

-- 2023-04-19: Add custom preferences

alter table users add column custom_preferences jsonb not null default '{}';
