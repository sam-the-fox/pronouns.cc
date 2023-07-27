-- +migrate Up

-- 2023-06-03: Add short IDs for the prns.cc domain.

-- add the columns
alter table users add column sid text unique check(length(sid)=5);
alter table members add column sid text unique check(length(sid)=6);
alter table users add column last_sid_reroll timestamptz not null default now() - '1 hour'::interval;

-- create the generate short ID functions
-- these are copied from PluralKit's HID functions:
-- https://github.com/PluralKit/PluralKit/blob/e4a2930bf353af9406e48934569677d7de6dd90d/PluralKit.Core/Database/Functions/functions.sql#L118-L152

-- +migrate StatementBegin
create function generate_sid(len int) returns text as $$
    select string_agg(substr('abcdefghijklmnopqrstuvwxyz', ceil(random() * 26)::integer, 1), '') from generate_series(1, len)
$$ language sql volatile;
-- +migrate StatementEnd

-- +migrate StatementBegin
create function find_free_user_sid() returns text as $$
declare new_sid text;
begin
    loop
        new_sid := generate_sid(5);
        if not exists (select 1 from users where sid = new_sid) then return new_sid; end if;
    end loop;
end
$$ language plpgsql volatile;
-- +migrate StatementEnd

-- +migrate StatementBegin
create function find_free_member_sid() returns text as $$
declare new_sid text;
begin
    loop
        new_sid := generate_sid(6);
        if not exists (select 1 from members where sid = new_sid) then return new_sid; end if;
    end loop;
end
$$ language plpgsql volatile;
-- +migrate StatementEnd

-- give all users and members short IDs
update users set sid = find_free_user_sid();
update members set sid = find_free_member_sid();

-- finally, make the values non-nullable
alter table users alter column sid set not null;
alter table members alter column sid set not null;
