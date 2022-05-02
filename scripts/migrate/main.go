// migrate runs (forward) migrations
package main

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	migrate "github.com/rubenv/sql-migrate"

	// SQL driver
	_ "github.com/jackc/pgx/v4/stdlib"
)

//go:embed *.sql
var migrations embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file:", err)
		os.Exit(1)
	}

	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("error opening database:", err)
		os.Exit(1)
	}

	fmt.Println("opened database")

	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("error closing database:", err)
			os.Exit(1)
		}
	}()

	if err := db.Ping(); err != nil {
		fmt.Println("error pinging database:", err)
		return
	}

	src := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrations,
		Root:       ".",
	}

	fmt.Println("executing migrations")
	n, err := migrate.Exec(db, "postgres", src, migrate.Up)
	if err != nil {
		fmt.Println("error executing migrations:", err)
		return
	}

	if n == 0 {
		fmt.Println("database is already up to date")
	} else {
		fmt.Printf("executed %v migrations!\n", n)
	}
}
