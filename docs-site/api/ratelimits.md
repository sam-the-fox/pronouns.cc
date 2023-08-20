# Rate limits

The API has rate limits, generally separated by groups of endpoints.
If you exceed a rate limit, the API will start to return 429 errors.

## Headers

- `X-RateLimit-Bucket`: the bucket the rate limit is for (listed below)
- `X-RateLimit-Limit`: the total number of requests you can make per minute
- `X-RateLimit-Remaining`: the number of requests remaining in the current timeframe
- `X-RateLimit-Reset`: the unix timestamp that the number of requests resets at
- `Retry-After`: only if you hit a rate limit, the number of seconds until you can make requests again

## Buckets

Note that only the most specific matching bucket is used for rate limits.

| Bucket                   | Rate limit per minute | Notes                                                       |
| ------------------------ | --------------------- | ----------------------------------------------------------- |
| /                        | 120                   | Used as fallback if no other bucket exists for the endpoint |
| GET /users/\*            | 60                    |                                                             |
| GET /users/\*/members    | 60                    |                                                             |
| GET /users/\*/members/\* | 60                    |                                                             |
| PATCH /users/@me         | 10                    |                                                             |
| POST /members            | 10                    |                                                             |
| GET /members/\*          | 60                    |                                                             |
| PATCH /members/\*        | 20                    |                                                             |
| DELETE /members/\*       | 5                     |                                                             |
| /auth/\*                 | 20                    |                                                             |
| /auth/tokens             | 10                    |                                                             |
| /auth/invites            | 10                    |                                                             |
