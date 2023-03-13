-- +migrate Up

-- 2023-03-13: Change avatar URLs to hashes

alter table users drop column avatar_urls;
alter table members drop column avatar_urls;

alter table users add column avatar text;
alter table members add column avatar text;
