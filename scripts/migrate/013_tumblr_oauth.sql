-- +migrate Up

-- 2023-04-18: Add tumblr oauth

alter table users add column tumblr text null;
alter table users add column tumblr_username text null;
