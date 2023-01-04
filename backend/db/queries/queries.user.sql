-- name: GetUserByID :one
SELECT * FROM users WHERE id = pggen.arg('id');

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = pggen.arg('username');

-- name: UpdateUserNamesPronouns :one
UPDATE users SET
names = pggen.arg('names'),
pronouns = pggen.arg('pronouns')
WHERE id = pggen.arg('id')
RETURNING *;

-- name: GetUserFields :many
SELECT * FROM user_fields WHERE user_id = pggen.arg('user_id') ORDER BY id ASC;

-- name: InsertUserField :one
INSERT INTO user_fields
(user_id, name, entries) VALUES
(pggen.arg('user_id'), pggen.arg('name'), pggen.arg('entries'))
RETURNING *;
