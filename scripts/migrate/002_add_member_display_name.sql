-- +migrate Up

-- 2022-11-20: add display name to members
alter table members add column display_name text;
