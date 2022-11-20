# pronouns.cc

A work-in-progress site to share your names, pronouns, and other preferred terms.

## Stack

- API server is written in Go with the [chi](https://github.com/go-chi/chi) router
- Persistent data is stored in PostgreSQL
- Temporary data is stored in Redis
- The frontend is written in TypeScript with React, using [Next](https://nextjs.org/) for server-side rendering
- Avatars are stored in S3-compatible storage ([MinIO](https://github.com/minio/minio) for development)

## Development

When working on the frontend, run the API and then use `yarn dev` in `frontend/` for hot reloading.

## Building

Run `make backend` to build the API server, then run `yarn build` in `frontend/`.

## Running

Both the backend and frontend are expected to run behind a reverse proxy such as [Caddy](https://caddyserver.com/).

The frontend should serve every possible path _except_ anything starting with `/api/`, which should be routed to the backend instead.
**Make sure to rewrite requests going to the API server to remove the starting `/api/`.**

## Development

Requirements:

- Go 1.18 or later
- PostgreSQL (any currently supported version should work)
- Redis 6.0 or later
- Node.js (latest version)
- MinIO **if testing uploads** (_not_ required otherwise)

### Setup

1. Create a PostgreSQL user and database (the user should own the database)  
   For example: `create user pronouns with password 'password'; create database pronouns with owner pronouns;`
2. Create a `.env` file in the repository root containing at least `HMAC_KEY`, `DATABASE_URL`, `REDIS`, `PORT`, and `MINIO_ENDPOINT` keys.
   (See below for an example)
3. Run `make migrate` to initialize the database, then optionally `make seeddb` to insert a test user.
4. Run `go run -v ./backend` to run the backend.
5. Create `frontend/.env` with the following content: `API_BASE=http://localhost:8080`
6. cd into the `frontend` directory and run `yarn dev` to run the frontend.

```sh
# Example env file
HMAC_KEY="`go run -v ./scripts/genkey`"
DATABASE_URL=postgresql://<username>:<pass>@localhost/<database> # PostgreSQL database URL
REDIS=localhost:6379
PORT=8080 # Port the API will listen on. Default is 8080, this is also default for the backend.
MINIO_ENDPOINT=localhost:9000 # This always needs to be set, it *does not* need to point to a MinIO server.
```

## License

    Copyright (C) 2022  Sam <u1f320>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published
    by the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
