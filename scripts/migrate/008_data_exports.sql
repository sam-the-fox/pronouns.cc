-- +migrate Up

-- 2023-03-15: Add data export

create table data_exports (
    id         serial primary key,
    user_id    text not null references users (id) on delete cascade,
    hash       text not null,
    created_at timestamptz not null default now()
);
