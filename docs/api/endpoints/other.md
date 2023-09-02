# Other endpoints

There are some endpoints that are neither user or member related:

### Get statistics

#### `GET /meta`

Get aggregate statistics for pronouns.cc.
Note: a user is considered active if they have updated their profile, created a member, deleted a member,
or updated a member's profile in the given time period.

#### Response body

| Field          | Type              | Description                                                               |
| -------------- | ----------------- | ------------------------------------------------------------------------- |
| git_repository | string            | link to the project's Git repository                                      |
| git_commit     | string            | the commit the backend is built from                                      |
| users          | user count object | the total number of users                                                 |
| members        | int               | the total number of non-hidden members                                    |
| require_invite | bool              | whether invites are required to sign up. _Always `false` for pronouns.cc_ |

#### User count object

| Field        | Type | Description                                 |
| ------------ | ---- | ------------------------------------------- |
| total        | int  | total number of users                       |
| active_month | int  | number of users active in the last month    |
| active_week  | int  | number of users active in the last week     |
| active_day   | int  | number of users active in the last 24 hours |

### Get warnings

#### `GET /auth/warnings`

**Requires authentication.** Returns an array of warnings the currently authenticated user has.
Add `?all=true` query parameter to return all warnings, not just unread ones.

#### Response body

| Field      | Type     | Description                                    |
| ---------- | -------- | ---------------------------------------------- |
| id         | int      | the warning ID                                 |
| reason     | string   | the reason for the warning                     |
| created_at | datetime | when the warning was created                   |
| read       | bool     | whether the warning has been read/acknowledged |
