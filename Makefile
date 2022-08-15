.PHONY: migrate
migrate:
	go run -v ./scripts/migrate

.PHONY: backend
backend:
	CGO_ENABLED=0 go build -v -o pronouns -ldflags="-buildid= -X codeberg.org/u1f320/pronouns.cc/backend/server.Revision=`git rev-parse --short HEAD`" ./backend
