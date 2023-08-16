-- 2022-11-20: add display name to members

-- +migrate Up

alter table members add column display_name text;

-- +migrate Down

alter table members drop column display_name;
