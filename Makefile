.PHONY: migrate
migrate:
	go run -v ./scripts/migrate

.PHONY: backend
backend:
	go build -v -o api -ldflags="-buildid= -X gitlab.com/1f320/pronouns/backend/server.Revision=`git rev-parse --short HEAD`" ./backend

.PHONY: frontend
frontend:
	yarn build

.PHONY: dev
dev:
	yarn dev
