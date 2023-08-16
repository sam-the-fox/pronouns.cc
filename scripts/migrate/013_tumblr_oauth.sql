-- 2023-04-18: Add tumblr oauth

-- +migrate Up

alter table users add column tumblr text null;
alter table users add column tumblr_username text null;

-- +migrate Down

alter table users drop column tumblr;
alter table users drop column tumblr_username;
