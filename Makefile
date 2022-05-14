.PHONY: migrate
migrate:
	go run -v ./scripts/migrate

.PHONY: backend
backend:
	go build -v -o api -ldflags="-buildid= -X codeberg.org/u1f320/pronouns.cc/backend/server.Revision=`git rev-parse --short HEAD`" ./backend

.PHONY: frontend
frontend:
	yarn build

.PHONY: dev
dev:
	yarn dev
