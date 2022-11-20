package main

import (
	"context"
	"fmt"
	"os"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file:", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("error opening database:", err)
		os.Exit(1)
	}
	defer pool.Close()

	fmt.Println("opened database")

	pg := &db.DB{Pool: pool}

	tx, err := pg.Begin(ctx)
	if err != nil {
		fmt.Println("error beginning transaction:", err)
		os.Exit(1)
	}

	u, err := pg.CreateUser(ctx, tx, "test")
	if err != nil {
		fmt.Println("error creating user:", err)
		os.Exit(1)
	}

	_, err = pg.UpdateUser(ctx, tx, u.ID, ptr("testing"), ptr("This is a bio!"), &[]string{"https://pronouns.cc"}, nil)
	if err != nil {
		fmt.Println("error setting user info:", err)
		os.Exit(1)
	}

	err = pg.SetUserNames(ctx, tx, u.ID, []db.Name{
		{Name: "testing 1", Status: db.StatusFavourite},
		{Name: "testing 2", Status: db.StatusOkay},
	})
	if err != nil {
		fmt.Println("error setting names:", err)
		os.Exit(1)
	}

	err = pg.SetUserPronouns(ctx, tx, u.ID, []db.Pronoun{
		{Pronouns: "it/it/its/its/itself", DisplayText: ptr("it/its"), Status: db.StatusFavourite},
		{Pronouns: "they/them/their/theirs/themself", Status: db.StatusOkay},
	})
	if err != nil {
		fmt.Println("error setting pronouns:", err)
		os.Exit(1)
	}

	err = pg.SetUserFields(ctx, tx, u.ID, []db.Field{
		{
			Name:        "Field 1",
			Favourite:   []string{"Favourite 1"},
			Okay:        []string{"Okay 1"},
			Jokingly:    []string{"Jokingly 1"},
			FriendsOnly: []string{"Friends only 1"},
			Avoid:       []string{"Avoid 1"},
		},
		{
			Name:        "Field 2",
			Favourite:   []string{"Favourite 2"},
			Okay:        []string{"Okay 2"},
			Jokingly:    []string{"Jokingly 2"},
			FriendsOnly: []string{"Friends only 2"},
			Avoid:       []string{"Avoid 2"},
		},
		{
			Name:        "Field 3",
			Favourite:   []string{"Favourite 3"},
			Okay:        []string{"Okay 3"},
			Jokingly:    []string{"Jokingly 3"},
			FriendsOnly: []string{"Friends only 3"},
			Avoid:       []string{"Avoid 3"},
		},
		{
			Name:        "Field 4",
			Favourite:   []string{"Favourite 4"},
			Okay:        []string{"Okay 4"},
			Jokingly:    []string{"Jokingly 4"},
			FriendsOnly: []string{"Friends only 4"},
			Avoid:       []string{"Avoid 4"},
		},
		{
			Name:        "Field 5",
			Favourite:   []string{"Favourite 5"},
			Okay:        []string{"Okay 5"},
			Jokingly:    []string{"Jokingly 5"},
			FriendsOnly: []string{"Friends only 5"},
			Avoid:       []string{"Avoid 5"},
		},
	})
	if err != nil {
		fmt.Println("error setting fields:", err)
		os.Exit(1)
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Println("error committing transaction:", err)
		os.Exit(1)
	}

	fmt.Println("Created testing user with ID", u.ID, "and name", u.Username)
}

func ptr[T any](v T) *T {
	return &v
}
