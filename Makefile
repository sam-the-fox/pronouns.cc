all: generate backend frontend

.PHONY: migrate
migrate:
	go run -v ./scripts/migrate

.PHONY: seeddb
seeddb:
	go run -v ./scripts/seeddb

.PHONY: backend
backend:
	CGO_ENABLED=0 go build -v -o pronouns -ldflags="-buildid= -X codeberg.org/u1f320/pronouns.cc/backend/server.Revision=`git rev-parse --short HEAD` -X codeberg.org/u1f320/pronouns.cc/backend/server.Tag=`git describe --tags --long`" .

.PHONY: generate
generate:
	go generate ./...

.PHONY: frontend
frontend:
	cd frontend && pnpm install && pnpm build

.PHONY: dev
dev:
	cd frontend && pnpm dev
