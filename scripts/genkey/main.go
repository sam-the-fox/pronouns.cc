package genkey

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "key",
	Usage:  "Generate a secure 64-byte base 64 key",
	Action: run,
}

func run(c *cli.Context) error {
	b := make([]byte, 64)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	s := base64.URLEncoding.EncodeToString(b)

	fmt.Println(s)
	return nil
}
