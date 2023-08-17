-- 2023-08-17: Add snowflake ID columns

-- +migrate Up

alter table users add column snowflake_id bigint unique;
alter table members add column snowflake_id bigint unique;
alter table pride_flags add column snowflake_id bigint unique;

-- +migrate Down

alter table users drop column snowflake_id;
alter table members drop column snowflake_id;
alter table pride_flags drop column snowflake_id;
