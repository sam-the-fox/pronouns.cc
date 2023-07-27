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

## Development

Requirements:

- Go 1.18 or later
- PostgreSQL (any currently supported version should work)
- Redis 6.0 or later
- Node.js (latest version)
- MinIO **if using avatars, flags, or data exports** (_not_ required otherwise)

### Setup

1. Create a PostgreSQL user and database (the user should own the database).  
   For example: `create user pronouns with password 'password'; create database pronouns with owner pronouns;`
2. Copy `.env.example` in the repository root to a new file named `.env` and fill out the required options.
3. Run `go run -v . database migrate` to initialize the database, then optionally `go run -v . database seed` to insert a test user.
4. Run `go run -v . web` to run the backend.
5. Copy `frontend/.env.example` into `frontend/.env` and tweak as necessary.
6. cd into the `frontend` directory and run `pnpm dev` to run the frontend.

See [`docs/production.md`](/docs/production.md#configuration) for more information about keys in the backend and frontend `.env` files.

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
