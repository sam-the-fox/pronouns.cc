-- +migrate Up

-- 2022-05-02: initial schema

create table users (
    id           text primary key,
    username     text not null unique,
    display_name text,
    bio          text,

    avatar_source text,
    avatar_url    text,
    links         text[],

    discord          text unique, -- for Discord oauth
    discord_username text
);

create table user_fields (
    user_id text      not null references users (id) on delete cascade,
    id      bigserial primary key,
    name    text      not null,

    favourite    text[] not null default array[]::text[],
    okay         text[] not null default array[]::text[],
    jokingly     text[] not null default array[]::text[],
    friends_only text[] not null default array[]::text[],
    avoid        text[] not null default array[]::text[]
);

create table members (
    id      text primary key,
    user_id text not null references users (id) on delete cascade,
    name    text not null,
    bio     text,

    avatar_url text,
    links      text
);

create table member_fields (
    member_id text      not null references members (id) on delete cascade,
    id        bigserial primary key,
    name      text      not null,

    favourite    text[] not null default array[]::text[],
    okay         text[] not null default array[]::text[],
    jokingly     text[] not null default array[]::text[],
    friends_only text[] not null default array[]::text[],
    avoid        text[] not null default array[]::text[]
);
