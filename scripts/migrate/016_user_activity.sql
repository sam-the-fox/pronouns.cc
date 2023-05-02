-- +migrate Up

-- 2023-05-02: Add a last_active column to users, updated whenever the user modifies their profile or members.
-- This is not directly exposed in the API.
-- Potential future use cases: showing total number of active users, pruning completely empty users if they don't log in?

alter table users add column last_active timestamptz not null default now();
