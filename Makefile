all: generate backend frontend

.PHONY: backend
backend:
	go build -v -o pronouns -ldflags="-buildid= -X codeberg.org/pronounscc/pronouns.cc/backend/server.Revision=`git rev-parse --short HEAD` -X codeberg.org/pronounscc/pronouns.cc/backend/server.Tag=`git describe --tags --long --always`" .

.PHONY: generate
generate:
	go generate ./...

.PHONY: frontend
frontend:
	cd frontend && pnpm install && pnpm build

.PHONY: dev
dev:
	cd frontend && pnpm dev
