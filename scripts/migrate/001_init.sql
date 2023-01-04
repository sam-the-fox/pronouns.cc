-- +migrate Up

-- 2022-05-02: initial schema

create table users (
    id           text primary key,
    username     text not null unique,
    display_name text,
    bio          text,

    avatar_urls text[],
    links       text[],

    discord          text unique, -- for Discord oauth
    discord_username text,

    max_invites int not null default 10
);

create table user_names (
    user_id text      not null references users (id) on delete cascade,
    id      bigserial primary key, -- ID is used for sorting; when order changes, existing rows are deleted and new ones are created
    name    text      not null,
    status  int       not null
);

create table user_pronouns (
    user_id      text      not null references users (id) on delete cascade,
    id           bigserial primary key,
    display_text text, -- if unset, falls back to first 2 parts of pronouns
    pronouns     text      not null,
    status       int       not null
);

create table user_fields (
    user_id text      not null references users (id) on delete cascade,
    id      bigserial primary key,
    name    text      not null,

    favourite    text[],
    okay         text[],
    jokingly     text[],
    friends_only text[],
    avoid        text[]
);

create table members (
    id      text primary key,
    user_id text not null references users (id) on delete cascade,
    name    text not null,
    bio     text,

    avatar_urls text[],
    links       text[]
);

create unique index members_user_name_idx on members (user_id, lower(name));

create table member_names (
    member_id text      not null references members (id) on delete cascade,
    id      bigserial primary key, -- ID is used for sorting; when order changes, existing rows are deleted and new ones are created
    name    text      not null,
    status  int       not null
);

create table member_pronouns (
    member_id    text      not null references members (id) on delete cascade,
    id           bigserial primary key,
    display_text text, -- if unset, falls back to first 2 parts of pronouns
    pronouns     text      not null,
    status       int       not null
);

create table member_fields (
    member_id text      not null references members (id) on delete cascade,
    id        bigserial primary key,
    name      text      not null,

    favourite    text[],
    okay         text[],
    jokingly     text[],
    friends_only text[],
    avoid        text[]
);

create table invites (
    user_id text      not null references users (id) on delete cascade,
    code    text      primary key,
    created timestamp not null default (current_timestamp at time zone 'utc'),
    used    boolean not null default false
);
