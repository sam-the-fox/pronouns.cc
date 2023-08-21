# User endpoints

## User object

| Field              | Type                                                 | Description                                                                 |
| ------------------ | ---------------------------------------------------- | --------------------------------------------------------------------------- |
| id                 | string                                               | the user's unique ID                                                        |
| sid                | string                                               | the user's 5 letter short ID                                                |
| name               | string                                               | the user's username                                                         |
| display_name       | string?                                              | the user's display name or nickname                                         |
| bio                | string?                                              | the user's description or bio                                               |
| member_title       | string?                                              | the heading used for the user's member list. If null, defaults to "Members" |
| avatar             | string?                                              | the user's [avatar hash](/api/#images)                                      |
| links              | string[]                                             | the user's profile links                                                    |
| names              | [field_entry](./#field-entry)[]                      | the user's preferred names                                                  |
| pronouns           | [pronoun_entry](./#pronoun-entry)[]                  | the user's preferred pronouns                                               |
| fields             | [field](./#field)[]                                  | the user's term fields                                                      |
| flags              | [flag](./#pride-flag)[]                              | the user's pride flags                                                      |
| members            | [partial](./users#partial-member-object) member[]    | the user's non-hidden members                                               |
| badges             | int                                                  | the user's badges, represented as a bitmask field                           |
| utc_offset         | int?                                                 | the user's current offset from UTC, in seconds                              |
| custom_preferences | map\[uuid\][custom_preference](./#custom-preference) | the user's custom preferences                                               |

### Additional fields for the currently authenticated user {#additional-user-fields}

| Field              | Type     | Description                                      |
| ------------------ | -------- | ------------------------------------------------ |
| created_at         | datetime | the user's creation date and time                |
| timezone           | string?  | the user's timezone in IANA timezone format      |
| is_admin           | bool     | whether or not the user is an administrator      |
| list_private       | bool     | whether or not the user's member list is private |
| last_sid_reroll    | datetime | the last time the user rerolled a short ID       |
| discord            | string?  | the user's Discord ID                            |
| discord_username   | string?  | the user's Discord username                      |
| tumblr             | string?  | the user's Tumblr ID                             |
| tumblr_username    | string?  | the user's Tumblr username                       |
| google             | string?  | the user's Google ID                             |
| google_username    | string?  | the user's Google username                       |
| fediverse          | string?  | the user's fediverse user ID                     |
| fediverse_username | string?  | the user's fediverse username, without instance  |
| fediverse_instance | string?  | the user's fediverse instance                    |

## Partial member object

| Field        | Type                                | Description                              |
| ------------ | ----------------------------------- | ---------------------------------------- |
| id           | string                              | the member's unique ID                   |
| sid          | string                              | the member's 6-letter short ID           |
| name         | string                              | the member's name                        |
| display_name | string?                             | the member's display name or nickname    |
| bio          | string?                             | the member's description                 |
| avatar       | string?                             | the member's [avatar hash](/api/#images) |
| links        | string[]                            | the member's profile links               |
| names        | [field_entry](./#field-entry)[]     | the member's preferred names             |
| pronouns     | [pronoun_entry](./#pronoun-entry)[] | the member's preferred pronouns          |

## Endpoints

### Get user

#### `GET /users/{user.id} | GET /users/{user.name}`

Gets a user by their ID or username. Returns a [user](./users#user-object) object.
If authenticated and the authenticated user is the requested user, also returns the [additional user fields](./users#additional-user-fields).

### Get current user

#### `GET /users/@me`

**Requires authentication.** Gets the currently authenticated [user](./users#user-object),
with all [additional user fields](./users/#additional-user-fields).

### Update current user

#### `PATCH /users/@me`

**Requires authentication.** Updates the currently authenticated user.
Returns the updated [user](./users#user-object) object on success.

#### Request body parameters

| Field              | Type                 | Description                                                                                                   |
| ------------------ | -------------------- | ------------------------------------------------------------------------------------------------------------- |
| name               | string               | the user's new name. Must be between 2 and 40 characters and only consist of ASCII letters, `_`, `.`, and `-` |
| display_name       | string               | the user's new display name. Must be between 1 and 100 characters                                             |
| bio                | string               | the user's new bio. Must be between 1 and 1000 characters                                                     |
| member_title       | string               | the user's new member title. Must be between 1 and 150 characters                                             |
| links              | string[]             | the user's new profile links. Maximum 25 links, and links must be between 1 and 256 characters                |
| names              | field_entry[]        | the users's new preferred names                                                                               |
| pronouns           | pronoun_entry[]      | the users's new preferred pronouns                                                                            |
| fields             | field[]              | the user's new profile fields                                                                                 |
| flags              | string[]             | the user's new flags. This must be an array of [pride flag](./#pride-flag) IDs.                               |
| avatar             | string               | the user's new avatar. This must be a PNG, JPEG, or WebP image, encoded in base64 data URI format             |
| timezone           | string               | the user's new timezone. Must be in IANA timezone database format                                             |
| list_private       | bool                 | whether or not the user's member list should be hidden                                                        |
| custom_preferences | _custom preferences_ | the user's new custom preferences                                                                             |

### Get pride flags

#### `GET /users/@me/flags`

**Requires authentication.** Returns an array of the currently authenticated user's [pride flags](./#pride-flag).

### Create pride flag

#### `POST /users/@me/flags`

**Requires authentication.** Creates a new pride flag. Returns a [pride flag](./#pride-flag) object on success.

#### Request body parameters

| Field       | Type   | Description                                                                                              |
| ----------- | ------ | -------------------------------------------------------------------------------------------------------- |
| flag        | string | the flag image. This must be a PNG, JPEG, or WebP image, encoded in base64 data URI format. **Required** |
| name        | string | the flag name. Must be between 1 and 100 characters. **Required**                                        |
| description | string | the flag description or alt text.                                                                        |

### Edit pride flag

#### `PATCH /users/@me/flags/{flag.id}`

**Requires authentication.** Edits an existing pride flag.
Returns the updated [pride flag](./#pride-flag) object on success.

#### Request body parameters

| Field       | Type   | Description                                                      |
| ----------- | ------ | ---------------------------------------------------------------- |
| name        | string | the flag's new name. Must be between 1 and 100 characters        |
| description | string | the flag's new description. Must be between 1 and 500 characters |

### Delete pride flag

#### `DELETE /users/@me/flags/{flag.id}`

**Requires authentication.** Deletes an existing pride flag. Returns `204 No Content` on success.

### Reroll short ID

#### `GET /users/@me/reroll`

**Requires authentication.** Rerolls the user's short ID. Returns the updated [user](./users#user-object) on success.
If the user has already rerolled a short ID in the past hour, returns `403 Forbidden`.
