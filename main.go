package main

import (
	"fmt"
	"os"

	"codeberg.org/u1f320/pronouns.cc/backend"
	"codeberg.org/u1f320/pronouns.cc/backend/exporter"
	"codeberg.org/u1f320/pronouns.cc/backend/prns"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"codeberg.org/u1f320/pronouns.cc/scripts/cleandb"
	"codeberg.org/u1f320/pronouns.cc/scripts/genid"
	"codeberg.org/u1f320/pronouns.cc/scripts/genkey"
	"codeberg.org/u1f320/pronouns.cc/scripts/migrate"
	"codeberg.org/u1f320/pronouns.cc/scripts/seeddb"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	HelpName: "pronouns.cc",
	Usage:    "Pronoun card website and API",
	Version:  server.Tag,
	Commands: []*cli.Command{
		backend.Command,
		exporter.Command,
		prns.Command,
		{
			Name:    "database",
			Aliases: []string{"db"},
			Usage:   "Manage the database",
			Subcommands: []*cli.Command{
				migrate.Command,
				seeddb.Command,
				cleandb.Command,
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"gen"},
			Usage:   "Generate various strings",
			Subcommands: []*cli.Command{
				genid.Command,
				genkey.Command,
			},
		},
	},
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
