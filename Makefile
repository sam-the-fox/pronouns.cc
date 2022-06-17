.PHONY: all
all: frontend css backend
	mv api pronouns

.PHONY: migrate
migrate:
	go run -v ./scripts/migrate

.PHONY: backend
backend: css
	CGO_ENABLED=0 go build -v -o api -ldflags="-buildid= -X codeberg.org/u1f320/pronouns.cc/backend/server.Revision=`git rev-parse --short HEAD`" ./backend

.PHONY: frontend
frontend:
	yarn build

.PHONY: css
css:
	yarn tailwindcss -m -o frontend/style.css

.PHONY: dev
dev:
	yarn dev
