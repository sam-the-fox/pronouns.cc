-- 2023-03-16: Add fediverse (Mastodon/Pleroma/Misskey) OAuth

-- +migrate Up

create table fediverse_apps (
    id            serial primary key,
    instance      text not null unique,
    client_id     text not null,
    client_secret text not null,
    instance_type text not null
);

alter table users add column fediverse text null;
alter table users add column fediverse_username text null;
alter table users add column fediverse_app_id integer null references fediverse_apps (id) on delete set null;

-- +migrate Down

alter table users drop column fediverse;
alter table users drop column fediverse_username;
alter table users drop column fediverse_app_id;

drop table fediverse_apps;
