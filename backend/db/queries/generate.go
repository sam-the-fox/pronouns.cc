package queries

//go:generate pggen gen go --query-glob queries.user.sql --query-glob queries.member.sql --postgres-connection "postgres://pggen:pggen@localhost/pggen"
