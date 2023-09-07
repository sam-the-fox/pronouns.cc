package main

import (
	"fmt"
	"os"

	"codeberg.org/pronounscc/pronouns.cc/backend"
	"codeberg.org/pronounscc/pronouns.cc/backend/exporter"
	"codeberg.org/pronounscc/pronouns.cc/backend/prns"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"codeberg.org/pronounscc/pronouns.cc/scripts/cleandb"
	"codeberg.org/pronounscc/pronouns.cc/scripts/genid"
	"codeberg.org/pronounscc/pronouns.cc/scripts/genkey"
	"codeberg.org/pronounscc/pronouns.cc/scripts/migrate"
	"codeberg.org/pronounscc/pronouns.cc/scripts/seeddb"
	"codeberg.org/pronounscc/pronouns.cc/scripts/snowflakes"
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
				snowflakes.Command,
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
