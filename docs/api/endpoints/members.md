# Member endpoints

## Member object

| Field        | Type                                                 | Description                                                                       |
| ------------ | ---------------------------------------------------- | --------------------------------------------------------------------------------- |
| id           | string                                               | the member's unique ID                                                            |
| id_new       | snowflake                                            | the member's unique snowflake ID                                                  |
| sid          | string                                               | the member's 6-letter short ID                                                    |
| name         | string                                               | the member's name                                                                 |
| display_name | string?                                              | the member's display name or nickname                                             |
| bio          | string?                                              | the member's description                                                          |
| avatar       | string?                                              | the member's [avatar hash](/api/#images)                                          |
| links        | string[]                                             | the member's profile links                                                        |
| names        | [field_entry](./#field-entry)[]                      | the member's preferred names                                                      |
| pronouns     | [pronoun_entry](./#pronoun-entry)[]                  | the member's preferred pronouns                                                   |
| fields       | ?[field](./#field)[]                                 | the member's term fields. Not returned in member list endpoints.                  |
| flags        | [flag](./#pride-flag)[]                              | the member's pride flags                                                          |
| user         | partial [user](./members#partial-user-object) object | the user associated with this member                                              |
| unlisted     | ?bool                                                | _only returned for your own members_, whether the member is shown in member lists |

## Partial user object

| Field              | Type                                                 | Description                            |
| ------------------ | ---------------------------------------------------- | -------------------------------------- |
| id                 | string                                               | the user's unique ID                   |
| id_new             | snowflake                                            | the user's unique snowflake ID         |
| name               | string                                               | the user's username                    |
| display_name       | string?                                              | the user's display name or nickname    |
| avatar             | string?                                              | the user's [avatar hash](/api/#images) |
| custom_preferences | map\[uuid\][custom_preference](./#custom-preference) | the user's custom preferences          |

## Endpoints

### Get member

#### `GET /members/{member.id}`

Gets a member by their ID. Returns a [member](./members#member-object) object.
If authenticated and the authenticated user is the owner of the requested member,
also returns the `unlisted` field.

### Get user member

#### `GET /users/{user.id}/members/{member.id} | GET /users/{user.name}/members/{member.name}`

Gets a member by their ID or name. Returns a [member](./members#member-object) object.
If authenticated and the authenticated user is the owner of the requested member,
also returns the `unlisted` field.

### Get user members

#### `GET /users/{user.id}/members | GET /users/{user.name}/members`

Get a user's members. Returns an array of [member](./members#member-object) objects.

### Get current user member

#### `GET /users/@me/members/{member.id} | GET /users/@me/members/{member.name}`

**Requires authentication.** Get one of the currently authenticated user's members by ID or name.
Returns a [member](./members#member-object) object.

### Get current user members

#### `GET /users/@me/members`

**Requires authentication.** Get the currently authenticated user's members.
Returns an array of [member](./members#member-object) objects.

### Create member

#### `POST /members`

**Requires authentication**. Creates a new member.
Returns the newly created [member](./members#member-object) on success.

#### Request body parameters

| Field        | Type            | Description                                                                                         |
| ------------ | --------------- | --------------------------------------------------------------------------------------------------- |
| name         | string          | the new member's name. Must be unique per user, and be between 1 and 100 characters. **Required**   |
| display_name | string?         | the new member's display name. Must be between 1 and 100 characters                                 |
| bio          | string?         | the new member's bio. Must be between 1 and 1000 characters                                         |
| avatar       | string          | the new member's avatar. This must be a PNG, JPEG, or WebP image, encoded in base64 data URI format |
| links        | string[]        | the new member's profile links. Maximum 25 links, and links must be between 1 and 256 characters    |
| names        | field_entry[]   | the new member's preferred names                                                                    |
| pronouns     | pronoun_entry[] | the new member's preferred pronouns                                                                 |
| fields       | field[]         | the new member's profile fields                                                                     |

### Update member

#### `PATCH /members/{member.id}`

**Requires authentication.** Updates the given member.
Returns the updated [member](./members#member-object) on success.

#### Request body parameters

| Field        | Type            | Description                                                                                            |
| ------------ | --------------- | ------------------------------------------------------------------------------------------------------ |
| name         | string          | the member's new name. Must be unique per user, and be between 1 and 100 characters.                   |
| display_name | string          | the member's new display name. Must be between 1 and 100 characters                                    |
| bio          | string          | the member's new bio. Must be between 1 and 1000 characters                                            |
| links        | string[]        | the member's new profile links. Maximum 25 links, and links must be between 1 and 256 characters       |
| names        | field_entry[]   | the member's new preferred names                                                                       |
| pronouns     | pronoun_entry[] | the member's new preferred pronouns                                                                    |
| fields       | field[]         | the member's new profile fields                                                                        |
| flags        | string[]        | the member's new flags. This must be an array of [pride flag](./#pride-flag) IDs, _not_ snowflake IDs. |
| avatar       | string          | the member's new avatar. This must be a PNG, JPEG, or WebP image, encoded in base64 data URI format    |
| unlisted     | bool            | whether or not the member should be hidden from the member list                                        |

### Delete member

#### `DELETE /members/{member.id}`

**Requires authentication.** Deletes the given member. Returns `204 No Content` on success.

### Reroll short ID

#### `GET /members/{member.id}/reroll`

**Requires authentication.** Rerolls the member's short ID.
Returns the updated [member](./members#member-object) on success.
If the user has already rerolled a short ID in the past hour, returns `403 Forbidden`.
