package seeddb

import (
	"fmt"
	"os"

	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "seed",
	Usage:  "Seed the database with test data",
	Action: run,
}

func run(c *cli.Context) error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file:", err)
		return err
	}

	ctx := c.Context

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("error opening database:", err)
		return err
	}
	defer pool.Close()

	fmt.Println("opened database")

	pg := &db.DB{Pool: pool}

	tx, err := pg.Begin(ctx)
	if err != nil {
		fmt.Println("error beginning transaction:", err)
		return err
	}

	u, err := pg.CreateUser(ctx, tx, "test")
	if err != nil {
		fmt.Println("error creating user:", err)
		return err
	}

	_, err = pg.UpdateUser(ctx, tx, u.ID, ptr("testing"), ptr("This is a bio!"), nil, ptr(false), &[]string{"https://pronouns.cc"}, nil, nil)
	if err != nil {
		fmt.Println("error setting user info:", err)
		return err
	}

	err = pg.SetUserNamesPronouns(ctx, tx, u.ID, []db.FieldEntry{
		{Value: "testing 1", Status: "favourite"},
		{Value: "testing 2", Status: "okay"},
	}, []db.PronounEntry{
		{Pronouns: "it/it/its/its/itself", DisplayText: ptr("it/its"), Status: "favourite"},
		{Pronouns: "they/them/their/theirs/themself", Status: "okay"},
	})
	if err != nil {
		fmt.Println("error setting pronouns:", err)
		return err
	}

	err = pg.SetUserFields(ctx, tx, u.ID, []db.Field{
		{
			Name: "Field 1",
			Entries: []db.FieldEntry{
				{Value: "Favourite 1", Status: "favourite"},
				{Value: "Okay 1", Status: "okay"},
				{Value: "Jokingly 1", Status: "jokingly"},
				{Value: "Friends only 1", Status: "friends_only"},
				{Value: "Avoid 1", Status: "avoid"},
			},
		},
		{
			Name: "Field 2",
			Entries: []db.FieldEntry{
				{Value: "Favourite 2", Status: "favourite"},
				{Value: "Okay 2", Status: "okay"},
				{Value: "Jokingly 2", Status: "jokingly"},
				{Value: "Friends only 2", Status: "friends_only"},
				{Value: "Avoid 2", Status: "avoid"},
			},
		},
		{
			Name: "Field 3",
			Entries: []db.FieldEntry{
				{Value: "Favourite 3", Status: "favourite"},
				{Value: "Okay 3", Status: "okay"},
				{Value: "Jokingly 3", Status: "jokingly"},
				{Value: "Friends only 3", Status: "friends_only"},
				{Value: "Avoid 3", Status: "avoid"},
			},
		},
		{
			Name: "Field 4",
			Entries: []db.FieldEntry{
				{Value: "Favourite 4", Status: "favourite"},
				{Value: "Okay 4", Status: "okay"},
				{Value: "Jokingly 4", Status: "jokingly"},
				{Value: "Friends only 4", Status: "friends_only"},
				{Value: "Avoid 4", Status: "avoid"},
			},
		},
		{
			Name: "Field 5",
			Entries: []db.FieldEntry{
				{Value: "Favourite 5", Status: "favourite"},
				{Value: "Okay 5", Status: "okay"},
				{Value: "Jokingly 5", Status: "jokingly"},
				{Value: "Friends only 5", Status: "friends_only"},
				{Value: "Avoid 5", Status: "avoid"},
			},
		},
	})
	if err != nil {
		fmt.Println("error setting fields:", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Println("error committing transaction:", err)
		return err
	}

	fmt.Println("Created testing user with ID", u.ID, "and name", u.Username)
	return nil
}

func ptr[T any](v T) *T {
	return &v
}
