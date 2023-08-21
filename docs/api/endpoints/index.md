# Object reference

These are some of the objects shared by multiple types of endpoints.
For other objects, such as [users](./users) or [members](./members), check their respective pages.

## Field

| Field   | Type                            | Description                 |
| ------- | ------------------------------- | --------------------------- |
| name    | string                          | the field's name or heading |
| entries | [field_entry](./#field-entry)[] | the field's entries         |

## Field entry

| Field  | Type   | Description                      |
| ------ | ------ | -------------------------------- |
| value  | string | this entry's value or key        |
| status | string | this entry's [status](./#status) |

## Pronoun entry

| Field        | Type    | Description                                                                                           |
| ------------ | ------- | ----------------------------------------------------------------------------------------------------- |
| pronouns     | string  | this entry's raw pronouns. This can be any user-inputted value and does not have to be a complete set |
| display_text | string? | the text shown in the pronoun list, if `pronouns` is a valid 5-member set                             |
| status       | string  | this entry's [status](./#status)                                                                      |

## Status

A name, pronoun, or field entry's **status** is how the user or member feels about that entry.
This can be any of `favourite`, `okay`, `jokingly`, `friends_only`, `avoid`,
as well as the UUID of any [custom preferences](./#custom-preference) the user has set.

## Custom preference

A user can set custom word preferences, which can have custom icons and tooltips. These are identified by a UUID.

| Field     | Type   | Description                                                                                          |
| --------- | ------ | ---------------------------------------------------------------------------------------------------- |
| icon      | string | the [Bootstrap icon](https://icons.getbootstrap.com/) associated with this preference                |
| tooltip   | string | the description shown in the tooltip on hover or tap                                                 |
| size      | string | the size at which any entry with this preference will be shown, can be `large`, `normal`, or `small` |
| muted     | bool   | whether the preference is shown in a muted grey colour                                               |
| favourite | bool   | whether the preference is treated the same as `favourite` when building embeds                       |

## Pride flag

| Field       | Type    | Description                           |
| ----------- | ------- | ------------------------------------- |
| id          | string  | the flag's unique ID                  |
| hash        | string  | the flag's [image hash](/api/#images) |
| name        | string  | the flag's name                       |
| description | string? | the flag's description or alt text    |
