migrate:
	go build -v ./scripts/migrate

api:
	go build -v -o api -ldflags="-buildid= -X gitlab.com/1f320/pronouns/backend/server.Revision=`git rev-parse --short HEAD`" ./backend
