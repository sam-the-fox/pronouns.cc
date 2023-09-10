package cleandb

import (
	"context"
	"fmt"
	"os"
	"time"

	dbpkg "codeberg.org/pronounscc/pronouns.cc/backend/db"
	"github.com/georgysavva/scany/v2/pgxscan"
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

	changeUnusedUsernames := os.Getenv("DB_CLEAN_CHANGE_UNUSED_USERNAMES") == "true"

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

	fmt.Println("deleting expired export files")
	var exports []dbpkg.DataExport
	err = pgxscan.Select(ctx, db, &exports, "SELECT * FROM data_exports WHERE created_at < $1", time.Now().Add(-dbpkg.KeepExportTime))
	if err != nil {
		fmt.Println("error getting to-be-deleted export files:", err)
		return err
	}

	for _, de := range exports {
		err = db.DeleteExport(ctx, de)
		if err != nil {
			fmt.Printf("error deleting export %v: %v\n", de.ID, err)
			continue
		}

		fmt.Println("deleted export", de.ID)
	}

	fmt.Printf("deleted %v expired exports\n", len(exports))

	if changeUnusedUsernames {
		fmt.Println("cleaning unused usernames")

		tx, err := db.Begin(ctx)
		if err != nil {
			fmt.Printf("error starting transaction: %v\n", err)
			return err
		}
		defer tx.Rollback(ctx)

		inactiveUsers, err := db.InactiveUsers(ctx, tx)
		if err != nil {
			fmt.Printf("getting inactive users: %v\n", err)
			return err
		}

		for _, u := range inactiveUsers {
			err = db.UpdateUsername(ctx, tx, u.ID, fmt.Sprintf("inactive-user-%v", u.SnowflakeID))
			if err != nil {
				fmt.Printf("changing username for user %v: %v\n", u.SnowflakeID, err)
				return err
			}
		}

		err = tx.Commit(ctx)
		if err != nil {
			fmt.Printf("committing transaction: %v\n", err)
			return err
		}

		fmt.Printf("changed usernames for %v inactive users\n", len(inactiveUsers))
	} else {
		fmt.Println("not cleaning unused usernames")
	}

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
		fmt.Println("there are no users pending deletion\nfinished cleaning database!")
		return nil
	}

	for _, u := range users {
		members, err := db.UserMembers(ctx, u.ID, true)
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
	fmt.Println("finished cleaning database!")
	return nil
}
