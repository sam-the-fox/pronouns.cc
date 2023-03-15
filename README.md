# pronouns.cc

A work-in-progress site to share your names, pronouns, and other preferred terms.

## Stack

- API server is written in Go with the [chi](https://github.com/go-chi/chi) router
- Persistent data is stored in PostgreSQL
- Temporary data is stored in Redis
- The frontend is written in TypeScript with Svelte, using [SvelteKit](https://kit.svelte.dev/) for server-side rendering
- Avatars are stored in S3-compatible storage ([MinIO](https://github.com/minio/minio) for development)

## Development

When working on the frontend, run the API and then use `pnpm dev` in `frontend/` for hot reloading.

Note that the Vite dev server assumes that the backend listens on `:8080` and MinIO listens on `:9000`.
If these ports differ on your development environment, you must edit `vite.config.ts`.

## Building

Run `make backend` to build the API server, then run `pnpm build` in `frontend/`.

## Running

Both the backend and frontend are expected to run behind a reverse proxy such as [Caddy](https://caddyserver.com/) or nginx.

Every path should be proxied to the frontend, except:

- `/api/`: this should be proxied to the backend, with the URL being rewritten to remove `/api`
  (for example, a request to `$DOMAIN/api/v1/users/@me` should be proxied to `localhost:8080/v1/users/@me`)
- `/media/`: this should be proxied to your object storage.
  Make sure to rewrite `/media` into your storage bucket's name.

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
3. Create a `.env` file in the `frontend/` directory containing `DOMAIN=http://localhost:5173`.
4. Run `make migrate` to initialize the database, then optionally `make seeddb` to insert a test user.
5. Run `go run -v ./backend` to run the backend.
6. Create `frontend/.env` with the following content: `API_BASE=http://localhost:8080`
7. cd into the `frontend` directory and run `pnpm dev` to run the frontend.

```sh
# Example env file
HMAC_KEY="`go run -v ./scripts/genkey`"
DATABASE_URL=postgresql://<username>:<pass>@localhost/<database> # PostgreSQL database URL
REDIS=localhost:6379
PORT=8080 # Port the API will listen on. Default is 8080, this is also default for the backend.
MINIO_ENDPOINT=localhost:9000 # This always needs to be set, it *does not* need to point to a running MinIO server.
```

## Updating in production

1. Build the backend with `make backend`
2. Build the frontend with `cd frontend && pnpm install && pnpm build`
3. Stop the servers: `systemctl stop pronouns-api pronouns-fe`
4. Run migrations: `./pronouns database migrate`
5. Start the servers: `systemctl start pronouns-api pronouns-fe`

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
