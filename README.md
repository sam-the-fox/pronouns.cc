# pronouns.cc

A work-in-progress site to share your pronouns and preferred terms.

## Stack

- API server is written in Go with the [chi](https://github.com/go-chi/chi) router
- Persistent data is stored in PostgreSQL
- Temporary data is stored in Redis
- The frontend is written in TypeScript with React, using [Next](https://nextjs.org/) for server-side rendering

## Development

When working on the frontend, run the API and then use `yarn dev` in `frontend/` for hot reloading.

## Building

Run `make backend` to build the API server, then run `yarn build` in `frontend/`.

## Running

Both the backend and frontend are expected to run behind a reverse proxy such as [Caddy](https://caddyserver.com/).

The frontend should serve every possible path _except_ anything starting with `/api/`, which should be routed to the backend instead.
**Make sure to rewrite requests going to the API server to remove the starting `/api/`.**

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
