package genid

import (
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:  "id",
	Usage: "Generate a time-based ID",
	Flags: []cli.Flag{
		&cli.TimestampFlag{
			Name:    "timestamp",
			Aliases: []string{"T"},
			Usage:   "The timestamp to generate an ID for (format: 2006-01-02T15:04:05)",
			Layout:  "2006-01-02T15:04:05",
		},
		&cli.UintFlag{
			Name:    "days-ago",
			Aliases: []string{"D"},
			Usage:   "The number of days ago to generate an ID for",
			Value:   0,
		},
	},
	Action: run,
}

func run(c *cli.Context) error {
	if t := c.Timestamp("timestamp"); t != nil {
		fmt.Printf("ID for %v: %v\n", t.Format(time.RFC1123), xid.NewWithTime(*t))
		return nil
	}
	if daysAgo := c.Uint("days-ago"); daysAgo != 0 {
		t := time.Now().Add(time.Duration(-daysAgo) * 24 * time.Hour)

		fmt.Printf("ID for %v days ago: %v\n", daysAgo, xid.NewWithTime(t))
		return nil
	}

	fmt.Printf("ID for now: %v\n", xid.New())
	return nil
}
