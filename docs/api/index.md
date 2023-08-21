# API reference

pronouns.cc has a HTTP REST API to query and edit profiles, available at `https://pronouns.cc/api`.

## Versioning

The API is versioned, and versions must be explicitly specified for all endpoints.
The current, and only, available version is **1**.
The version is specified in the request path, like `https://pronouns.cc/api/v{version}`.

| Version | Status     |
| ------- | ---------- |
| 1       | Default    |
| 2       | _Upcoming_ |

The API version will be incremented for any breaking changes, including:

- Removing entire endpoints
- Removing fields from responses
- Changing the behaviour of fields (in some situations, see below)

However, the following types of changes are **not** considered breaking:

- Adding new endpoints
- Adding new fields to requests or responses (your JSON serializer/deserializer should ignore unknown fields)
- Forcing fields related to removed features to their default value

## Authentication

Tokens can be created [here](https://pronouns.cc/settings/tokens).
Not all endpoints require authentication. For those that do, a token must be provided in the `Authorization` header.
The token _may_ be prefixed with `Bearer `, but this is not required.

::: info
You are allowed to use site tokens (those stored in your browser's local storage) to access endpoints not available to API tokens,
however, these endpoints are not available to API tokens *for a reason*:
site tokens can take destructive actions such as deleting your account.
Additionally, endpoints that are not available to API tokens may have breaking changes without a major version bump.
:::

## Request bodies

::: info
The current API version doesn't distinguish between omitted and `null` keys yet.
However, the next version of the API will use `null` to unset keys, so clients should not rely on this behaviour.
:::

Request bodies should be in JSON format.
For PATCH requests, **all keys are optional**. Omitted keys will not be updated,
and keys set to the zero value of their respective types (for strings: `""`, for numbers: `0`, for arrays: `[]`, etc.)
will be unset.

## Documentation formatting

The "type" column in tables is formatted as follows:

- The type used is the _Go_ type, not the _JSON_ type.
  For example, the documentation will use `int` for integers and `float` for floats,
  even though they are both represented with JSON numbers.
- A _leading_ `?` signifies that the field may be omitted.
- A _trailing_ `?` signifies that the field may be null.

## IDs

::: info
pronouns.cc is [planning a transition](https://codeberg.org/pronounscc/pronouns.cc/issues/89)
to [Snowflake IDs](https://en.wikipedia.org/wiki/Snowflake_ID).
The information below pertains to the current ID format.
:::

The API uses [xid](https://github.com/rs/xid) for unique IDs. These are always serialized as strings.
Although xids have timestamp information embedded in them, this is non-trivial to extract.
xids are unique across _all_ resources, they are never shared (for example, a user and a member cannot share the same ID).

Users and members also have an additional ID type, `sid`.
These are randomly generated 5 or 6 letter strings, and are used for the prns.cc URL shortener.
They can be rerolled once per hour.

## Images

The API does not return full URLs to images such as avatars and pride flags.
Instead, the URL must be constructed manually using the `avatar` or `hash` fields.

The default user and member avatar is served at `https://pronouns.cc/default/512.webp`.
All custom images are served on the base URL `https://cdn.pronouns.cc`, and are only available in WebP format.

| Type          | Format                                      |
| ------------- | ------------------------------------------- |
| User avatar   | `/users/{user.id}/{user.avatar}.webp`       |
| Member avatar | `/members/{member.id}/{member.avatar}.webp` |
| Pride flag    | `/flags/{flag.hash}.webp`                   |
