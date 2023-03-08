-- +migrate Up

-- 2023-03-07: add delete functionality

-- if not null, the user is soft deleted
alter table users add column deleted_at timestamptz;
-- if true, the user deleted their account themselves + should have option to reactivate; should also be deleted after 30 days
alter table users add column self_delete boolean;
-- delete reason if the user was deleted by a moderator
alter table users add column delete_reason text;
