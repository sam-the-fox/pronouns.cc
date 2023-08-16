-- 2023-03-13: Change avatar URLs to hashes

-- +migrate Up

alter table users drop column avatar_urls;
alter table members drop column avatar_urls;

alter table users add column avatar text;
alter table members add column avatar text;

-- +migrate Down

alter table users drop column avatar;
alter table members drop column avatar;

alter table users add column avatar_urls text[];
alter table members add column avatar_urls text[];
