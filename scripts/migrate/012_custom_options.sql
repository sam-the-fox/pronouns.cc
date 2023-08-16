-- 2023-04-01: Add a couple customization options to users and members

-- +migrate Up

alter table users add column member_title text;
alter table users add column list_private boolean not null default false;

alter table members add column unlisted boolean not null default false;

-- +migrate Down

alter table users drop column member_title;
alter table users drop column list_private;

alter table members drop column unlisted;
