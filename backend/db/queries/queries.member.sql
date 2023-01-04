-- name: GetMemberByID :one
SELECT * FROM members
WHERE id = pggen.arg('id');

-- name: GetMemberByName :one
SELECT * FROM members
WHERE user_id = pggen.arg('user_id') AND (
    id = pggen.arg('member_ref')
    OR name = pggen.arg('member_ref')
);

-- name: GetMembers :many
SELECT * FROM members
WHERE user_id = pggen.arg('user_id')
ORDER BY name, id;

-- name: UpdateMemberNamesPronouns :one
UPDATE members SET
names = pggen.arg('names'),
pronouns = pggen.arg('pronouns')
WHERE id = pggen.arg('id')
RETURNING *;

-- name: GetMemberFields :many
SELECT * FROM member_fields WHERE member_id = pggen.arg('member_id') ORDER BY id ASC;

-- name: InsertMemberField :one
INSERT INTO member_fields
(member_id, name, entries) VALUES
(pggen.arg('member_id'), pggen.arg('name'), pggen.arg('entries'))
RETURNING *;
