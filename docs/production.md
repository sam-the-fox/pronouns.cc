# Running pronouns.cc in production

The configuration files in this directory are the same files used to run pronouns.cc in production.  
You might have to change paths and ports, but they should work fine as-is.

## Building pronouns.cc

```bash
git clone https://codeberg.org/u1f320/pronouns.cc.git pronouns
cd pronouns
make all

# if required fonts have not been downloaded yet
./download-fonts.sh
# if running for the first time
./pronouns database migrate
```

## Configuration

pronouns.cc is configured using `.env` files. Note that there are two separate `.env` files,
one in the repository root (for the backend) and one in the frontend directory.

### Backend keys

- `HMAC_KEY`: the key used to sign tokens. This should be a base64 string, you can generate one with `scripts/genkey`.
- `DATABASE_URL`: the URL for the PostgreSQL database.
- `REDIS`: the URL for the Redis database.
- `PORT` (int): the port the backend will listen on.
- `EXPORTER_PORT` (int): the port that the exporter service will listen on.
- `DEBUG` (true/false): whether to enable request logging.
- `BASE_URL`: the base URL for the frontend, used to construct some links.
- `MINIO_ENDPOINT`: the S3 endpoint for object storage.
- `MINIO_BUCKET`: the S3 bucket name.
- `MINIO_ACCESS_KEY_ID`: the S3 access key ID.
- `MINIO_ACCESS_KEY_SECRET`: the S3 access key secret.
- `MINIO_SSL`: whether to use SSL for S3.
- `FRONTEND_IP`: the IP for the frontend, which the rate limiter will ignore.
- `REQUIRE_INVITE`: whether to require invites to sign up.
- `DISCORD_CLIENT_ID`: for Discord auth, the client ID.
- `DISCORD_CLIENT_SECRET`: for Discord auth, the client secret.
- `DISCORD_PUBLIC_KEY`: public key for the Discord bot endpoint.

### Frontend keys

- `PUBLIC_BASE_URL`: the base URL for the frontend.
- `PRIVATE_SENTRY_DSN`: your Sentry DSN.

## Updating

```bash
make all
systemctl stop pronouns-api pronouns-fe
systemctl stop pronouns-exporter # only if the User, Member, Field, Export tables changed
./pronouns database migrate # only if a new migration was added
systemctl start pronouns-api pronouns-fe
systemctl start pronouns-exporter # if the exporter was stopped
```

## Proxy

Both the backend and frontend are expected to run behind a reverse proxy such as Caddy or nginx.
This directory contains a sample configuration file for nginx.

Every path should be proxied to the frontend, except:

- `/api/`: this should be proxied to the backend, with the URL being rewritten to remove `/api`
  (for example, a request to `$DOMAIN/api/v1/users/@me` should be proxied to `localhost:8080/v1/users/@me`)
- `/media/`: this should be proxied to your object storage.
  Make sure to rewrite `/media` into your storage bucket's name.