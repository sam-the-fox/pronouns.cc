# pronouns.cc

A work-in-progress site to share your pronouns and preferred terms.

## Stack

- API server is written in Go with the [chi](https://github.com/go-chi/chi) router
- Persistent data is stored in PostgreSQL
- Temporary data is stored in Redis
- The frontend is written in TypeScript with React, using [Vite](https://vitejs.dev/)

## Development

When working on the frontend, run the API and then use `yarn dev` for hot reloading.

Note that Vite does _not_ proxy the preview pages (`/@user` and `/@user/member`),
and as such, the Go server must be used directly to test it. However, as the preview pages are
rendered by the Go server, and doesn't share code with the frontend (except for CSS), I consider this an acceptable compromise.

## Building

Run `make all`. This will build the frontend, then embed that in the backend.

The resulting `pronouns` binary is a statically linked executable containing everything needed to run the website.  
Note that it should still be run behind a reverse proxy for TLS.

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
