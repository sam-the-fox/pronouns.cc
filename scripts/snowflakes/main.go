package snowflakes

import (
	"errors"
	"os"
	"time"

	"codeberg.org/pronounscc/pronouns.cc/backend/common"
	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/rs/xid"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "create-snowflakes",
	Usage:  "Give all users, members, and flags snowflake IDs.",
	Action: run,
}

func run(c *cli.Context) error {
	err := godotenv.Load()
	if err != nil {
		log.Error("loading .env file:", err)
		return err
	}

	conn, err := pgx.Connect(c.Context, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error("opening database:", err)
		return err
	}
	defer conn.Close(c.Context)
	log.Info("opened database")

	tx, err := conn.Begin(c.Context)
	if err != nil {
		log.Error("creating transaction:", err)
		return err
	}
	defer func() {
		err := tx.Rollback(c.Context)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("rolling back transaction:", err)
		}
	}()

	var userIDs []xid.ID
	err = pgxscan.Select(c.Context, conn, &userIDs, "SELECT id FROM users WHERE snowflake_id IS NULL")
	if err != nil {
		log.Error("selecting users without snowflake:", err)
		return err
	}

	t := time.Now()
	for _, userID := range userIDs {
		t := userID.Time()
		snowflake := common.UserID(common.GenerateIDWithTime(t))

		_, err = tx.Exec(c.Context, "UPDATE users SET snowflake_id = $1 WHERE id = $2", snowflake, userID)
		if err != nil {
			log.Errorf("updating user with ID %v: %v", userID, err)
			return err
		}
	}
	log.Infof("updated %v users in %v", len(userIDs), time.Since(t))

	var memberIDs []xid.ID
	err = pgxscan.Select(c.Context, conn, &memberIDs, "SELECT id FROM members WHERE snowflake_id IS NULL")
	if err != nil {
		log.Error("selecting users without snowflake:", err)
		return err
	}

	t = time.Now()
	for _, memberID := range memberIDs {
		t := memberID.Time()
		snowflake := common.MemberID(common.GenerateIDWithTime(t))

		_, err = tx.Exec(c.Context, "UPDATE members SET snowflake_id = $1 WHERE id = $2", snowflake, memberID)
		if err != nil {
			log.Errorf("updating user with ID %v: %v", memberID, err)
			return err
		}
	}
	log.Infof("updated %v members in %v", len(memberIDs), time.Since(t))

	var flagIDs []xid.ID
	err = pgxscan.Select(c.Context, conn, &flagIDs, "SELECT id FROM pride_flags WHERE snowflake_id IS NULL")
	if err != nil {
		log.Error("selecting users without snowflake:", err)
		return err
	}

	t = time.Now()
	for _, flagID := range flagIDs {
		t := flagID.Time()
		snowflake := common.FlagID(common.GenerateIDWithTime(t))

		_, err = tx.Exec(c.Context, "UPDATE pride_flags SET snowflake_id = $1 WHERE id = $2", snowflake, flagID)
		if err != nil {
			log.Errorf("updating user with ID %v: %v", flagID, err)
			return err
		}
	}
	log.Infof("updated %v flags in %v", len(flagIDs), time.Since(t))

	err = tx.Commit(c.Context)
	if err != nil {
		log.Error("committing transaction:", err)
		return err
	}

	return nil
}
