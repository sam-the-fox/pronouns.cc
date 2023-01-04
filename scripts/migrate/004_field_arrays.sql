-- +migrate Up

-- 2023-01-03: change names, pronouns, and fields to be columns instead of separate tables

create type field_entry as (
    value  text,
    status int
);

create type pronoun_entry as (
    value         text,
    display_value text,
    status        int
);

alter table users add column names field_entry[];
alter table users add column pronouns pronoun_entry[];

alter table members add column names field_entry[];
alter table members add column pronouns pronoun_entry[];

alter table user_fields add column entries field_entry[];
alter table member_fields add column entries field_entry[];

alter table user_fields drop column favourite;
alter table user_fields drop column okay;
alter table user_fields drop column jokingly;
alter table user_fields drop column friends_only;
alter table user_fields drop column avoid;

alter table member_fields drop column favourite;
alter table member_fields drop column okay;
alter table member_fields drop column jokingly;
alter table member_fields drop column friends_only;
alter table member_fields drop column avoid;
