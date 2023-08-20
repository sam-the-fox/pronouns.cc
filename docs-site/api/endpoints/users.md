# User endpoints

## User model

| Field        | Type    | Description                                                                 |
| ------------ | ------- | --------------------------------------------------------------------------- |
| id           | string  | the user's unique ID                                                        |
| sid          | string  | the user's 5 letter short ID                                                |
| name         | string  | the user's username                                                         |
| display_name | string? | the user's display name or nickname                                         |
| bio          | string? | the user's description or bio                                               |
| member_title | string? | the heading used for the user's member list. If null, defaults to "Members" |
| avatar       | string? | the user's [avatar hash](/api/#images)                                      |
