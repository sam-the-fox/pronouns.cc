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

	err = pg.SetUserNamesPronouns(ctx, tx, u.ID, []db.FieldEntry{
		{Value: "testing 1", Status: db.StatusFavourite},
		{Value: "testing 2", Status: db.StatusOkay},
	}, []db.PronounEntry{
		{Pronouns: "it/it/its/its/itself", DisplayText: ptr("it/its"), Status: db.StatusFavourite},
		{Pronouns: "they/them/their/theirs/themself", Status: db.StatusOkay},
	})
	if err != nil {
		fmt.Println("error setting pronouns:", err)
		os.Exit(1)
	}

	err = pg.SetUserFields(ctx, tx, u.ID, []db.Field{
		{
			Name: "Field 1",
			Entries: []db.FieldEntry{
				{Value: "Favourite 1", Status: db.StatusFavourite},
				{Value: "Okay 1", Status: db.StatusOkay},
				{Value: "Jokingly 1", Status: db.StatusJokingly},
				{Value: "Friends only 1", Status: db.StatusFriendsOnly},
				{Value: "Avoid 1", Status: db.StatusAvoid},
			},
		},
		{
			Name: "Field 2",
			Entries: []db.FieldEntry{
				{Value: "Favourite 2", Status: db.StatusFavourite},
				{Value: "Okay 2", Status: db.StatusOkay},
				{Value: "Jokingly 2", Status: db.StatusJokingly},
				{Value: "Friends only 2", Status: db.StatusFriendsOnly},
				{Value: "Avoid 2", Status: db.StatusAvoid},
			},
		},
		{
			Name: "Field 3",
			Entries: []db.FieldEntry{
				{Value: "Favourite 3", Status: db.StatusFavourite},
				{Value: "Okay 3", Status: db.StatusOkay},
				{Value: "Jokingly 3", Status: db.StatusJokingly},
				{Value: "Friends only 3", Status: db.StatusFriendsOnly},
				{Value: "Avoid 3", Status: db.StatusAvoid},
			},
		},
		{
			Name: "Field 4",
			Entries: []db.FieldEntry{
				{Value: "Favourite 4", Status: db.StatusFavourite},
				{Value: "Okay 4", Status: db.StatusOkay},
				{Value: "Jokingly 4", Status: db.StatusJokingly},
				{Value: "Friends only 4", Status: db.StatusFriendsOnly},
				{Value: "Avoid 4", Status: db.StatusAvoid},
			},
		},
		{
			Name: "Field 5",
			Entries: []db.FieldEntry{
				{Value: "Favourite 5", Status: db.StatusFavourite},
				{Value: "Okay 5", Status: db.StatusOkay},
				{Value: "Jokingly 5", Status: db.StatusJokingly},
				{Value: "Friends only 5", Status: db.StatusFriendsOnly},
				{Value: "Avoid 5", Status: db.StatusAvoid},
			},
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
