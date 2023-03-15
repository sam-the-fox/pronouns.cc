package cleandb

import (
	"context"
	"fmt"
	"time"

	dbpkg "codeberg.org/u1f320/pronouns.cc/backend/db"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/joho/godotenv"
	"github.com/rs/xid"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "clean",
	Usage:  "Clean deleted tokens + users daily",
	Action: run,
}

func run(c *cli.Context) error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file:", err)
		return err
	}

	ctx := context.Background()

	db, err := dbpkg.New()
	if err != nil {
		fmt.Println("error opening database:", err)
		return err
	}
	defer db.Close()

	fmt.Println("opened database")

	fmt.Println("deleting invalidated tokens")

	ct, err := db.Exec(ctx, "DELETE FROM tokens WHERE invalidated = true OR expires < $1", time.Now())
	if err != nil {
		fmt.Println("executing query:", err)
		return err
	}

	fmt.Printf("deleted %v invalidated or expired tokens\n", ct.RowsAffected())

	var users []dbpkg.User
	err = pgxscan.Select(ctx, db, &users, `SELECT * FROM users WHERE
	deleted_at IS NOT NULL AND
	(self_delete = true AND deleted_at < $1)
	OR (self_delete = false AND deleted_at < $2)
	ORDER BY id`, time.Now().Add(-dbpkg.SelfDeleteAfter), time.Now().Add(-dbpkg.ModDeleteAfter))
	if err != nil {
		fmt.Println("error getting to-be-deleted users:", err)
		return err
	}

	if len(users) == 0 {
		fmt.Println("there are no users pending deletion")
		return nil
	}

	for _, u := range users {
		members, err := db.UserMembers(ctx, u.ID)
		if err != nil {
			fmt.Printf("error getting members for user %v: %v\n", u.ID, err)
			continue
		}

		for _, m := range members {
			if m.Avatar == nil {
				continue
			}

			fmt.Printf("deleting avatars for member %v\n", m.ID)

			err = db.DeleteMemberAvatar(ctx, m.ID, *m.Avatar)
			if err != nil {
				fmt.Printf("error deleting avatars for member %v: %v", m.ID, err)
				continue
			}

			fmt.Printf("deleted avatars for member %v\n", m.ID)
		}

		if u.Avatar == nil {
			continue
		}

		fmt.Printf("deleting avatars for user %v\n", u.ID)

		err = db.DeleteUserAvatar(ctx, u.ID, *u.Avatar)
		if err != nil {
			fmt.Printf("error deleting avatars for user %v: %v", u.ID, err)
			continue
		}

		fmt.Printf("deleted avatars for user %v\n", u.ID)
	}

	ids := make([]xid.ID, len(users))
	for _, u := range users {
		ids = append(ids, u.ID)
	}

	ct, err = db.Exec(ctx, "DELETE FROM users WHERE id = ANY($1)", ids)
	if err != nil {
		fmt.Printf("error deleting users: %v\n", err)
		return err
	}

	fmt.Printf("deleted %v users!\n", ct.RowsAffected())
	return nil
}
