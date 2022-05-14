package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	port := ":" + os.Getenv("PORT")

	s, err := server.New()
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	// mount api routes
	mountRoutes(s)

	e := make(chan error)

	// run server in another goroutine (for gracefully shutting down, see below)
	go func() {
		e <- http.ListenAndServe(port, s.Router)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Infof("API server running at %v!", port)

	select {
	case <-ctx.Done():
		log.Info("Interrupt signal received, shutting down...")
		s.DB.Close()
		return
	case err := <-e:
		log.Fatalf("Error running server: %v", err)
	}
}
