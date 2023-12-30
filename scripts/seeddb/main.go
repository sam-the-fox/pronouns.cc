package seeddb

import (
	"errors"
	"log"
	"os"

	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type Seed struct {
	Users []SeedUser `yaml:"users"`
}

type SeedUser struct {
	Username    string            `yaml:"username"`
	DisplayName *string           `yaml:"displayName"`
	Bio         *string           `yaml:"bio"`
	Links       []string          `yaml:"links"`
	Names       []db.FieldEntry   `yaml:"names"`
	Pronouns    []db.PronounEntry `yaml:"pronouns"`
	Fields      []db.Field        `yaml:"fields"`
	Members     []SeedMember      `yaml:"members"`
}

type SeedMember struct {
	Name        string            `yaml:"name"`
	DisplayName *string           `yaml:"displayName"`
	Bio         string            `yaml:"bio"`
	Links       []string          `yaml:"links"`
	Names       []db.FieldEntry   `yaml:"names"`
	Pronouns    []db.PronounEntry `yaml:"pronouns"`
	Fields      []db.Field        `yaml:"fields"`
	Members     []SeedMember      `yaml:"members"`
}

var Command = &cli.Command{
	Name:   "seed",
	Usage:  "Seed the database with test data",
	Action: run,
}

func run(c *cli.Context) error {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file:", err)
		return err
	}

	ctx := c.Context

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("error opening database:", err)
		return err
	}
	defer pool.Close()

	log.Println("opened database")

	pg := &db.DB{Pool: pool}

	// read seed file
	seedFile, err := os.ReadFile("seed.yaml")
	if err != nil {
		log.Println("error opening seed.yaml:", err)
		return err
	}

	var seed Seed
	err = yaml.Unmarshal(seedFile, &seed)
	if err != nil {
		log.Println("error reading seed.yaml:", err)
		return err
	}

	tx, err := pg.Begin(ctx)
	if err != nil {
		log.Println("error beginning transaction:", err)
		return err
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Println("error rolling back transaction:", err)
		}
	}()

	for i, su := range seed.Users {
		u, err := pg.CreateUser(ctx, tx, su.Username)
		if err != nil {
			log.Printf("error creating user #%v/%s: %v", i+1, su.Username, err)
			return err
		}

		_, err = pg.UpdateUser(ctx, tx, u.ID, su.DisplayName, su.Bio, nil, nil, &su.Links, nil, nil, nil)
		if err != nil {
			log.Printf("updating user %s: %v", su.Username, err)
			return err
		}

		err = pg.SetUserNamesPronouns(ctx, tx, u.ID, db.NotNull(su.Names), db.NotNull(su.Pronouns))
		if err != nil {
			log.Printf("setting names/pronouns for user %s: %v", su.Username, err)
			return err
		}

		err = pg.SetUserFields(ctx, tx, u.ID, db.NotNull(su.Fields))
		if err != nil {
			log.Printf("setting fields for user %s: %v", su.Username, err)
			return err
		}

		log.Printf("creating members for user %s", su.Username)

		for _, sm := range su.Members {
			m, err := pg.CreateMember(ctx, tx, u.ID, sm.Name, sm.DisplayName, sm.Bio, db.NotNull(sm.Links))
			if err != nil {
				log.Printf("creating member %s: %v", sm.Name, err)
				return err
			}

			err = pg.SetMemberNamesPronouns(ctx, tx, m.ID, db.NotNull(sm.Names), db.NotNull(sm.Pronouns))
			if err != nil {
				log.Printf("setting names/pronouns for member %s: %v", sm.Name, err)
				return err
			}

			err = pg.SetMemberFields(ctx, tx, m.ID, db.NotNull(sm.Fields))
			if err != nil {
				log.Printf("setting fields for member %s: %v", sm.Name, err)
				return err
			}

			log.Printf("created member %s", sm.Name)
		}

		log.Printf("created user %s", su.Username)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("error committing transaction:", err)
		return err
	}

	log.Printf("seeded database with %d users", len(seed.Users))
	return nil
}
