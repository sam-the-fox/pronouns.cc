-- +migrate Up

-- 2023-04-01: Add a couple customization options to users and members

alter table users add column member_title text;
alter table users add column list_private boolean not null default false;

alter table members add column unlisted boolean not null default false;
