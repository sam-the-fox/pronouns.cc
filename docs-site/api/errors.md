# Error messages

If there is an error in your request, or the server encounters an error while processing it, an error object will be returned.

| Field           | Type    | Description                                                     |
| --------------- | ------- | --------------------------------------------------------------- |
| code            | int     | an [error code](./errors#error-codes)                           |
| message         | ?string | a human-readable description of the error                       |
| details         | ?string | more details about the error, most often for bad request errors |
| ratelimit_reset | ?int    | the unix time when an expired rate limit will reset             |

### Error codes

| Code | Description                                                                         |
| ---- | ----------------------------------------------------------------------------------- |
| 400  | One or more fields in your requests was invalid, or some required field is missing. |
| 403  | You are not authorized to use this endpoint.                                        |
| 404  | The endpoint was not found.                                                         |
| 405  | The method you are trying to use is not suported for this endpoint.                 |
| 429  | You have made too many requests in the last minute.                                 |
| 500  | An internal server error occurred.                                                  |
| 1006 | That username is invalid.                                                           |
| 1007 | That username is already taken.                                                     |
| 2001 | User not found.                                                                     |
| 2002 | This user's member list is private.                                                 |
| 2003 | You have reached the maximum number of pride flags.                                 |
| 2004 | You are trying to reroll short IDs too quickly.                                     |
| 3001 | Member not found.                                                                   |
| 3002 | You have reached the maximum number of members.                                     |
| 3003 | That member name is already in use.                                                 |
| 3004 | You can only edit your own members.                                                 |
| 4001 | Your request is too big (maximum 2 megabytes)                                       |
| 4002 | This endpoint is unavailable to your account or the current token.                  |
