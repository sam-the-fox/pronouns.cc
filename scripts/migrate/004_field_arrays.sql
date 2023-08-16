-- 2023-01-03: change names, pronouns, and fields to be columns instead of separate tables

-- +migrate Up

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

-- +migrate Down

alter table user_fields add column favourite text[];
alter table user_fields add column okay text[];
alter table user_fields add column jokingly text[];
alter table user_fields add column friends_only text[];
alter table user_fields add column avoid text[];

alter table member_fields add column favourite text[];
alter table member_fields add column okay text[];
alter table member_fields add column jokingly text[];
alter table member_fields add column friends_only text[];
alter table member_fields add column avoid text[];

alter table users drop column names;
alter table users drop column pronouns;

alter table members drop column names;
alter table members drop column pronouns;

alter table user_fields drop column entries;
alter table member_fields drop column entries;

drop type field_entry;
drop type pronoun_entry;
