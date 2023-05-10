package backend

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"

	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/go-chi/render"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "web",
	Usage:  "Run the API server",
	Action: run,
}

func run(c *cli.Context) error {
	vips.Startup(nil)
	defer vips.Shutdown()

	port := ":" + os.Getenv("PORT")

	s, err := server.New()
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
		return nil
	}

	// set render.Decode to a custom one that checks content length
	render.Decode = decode

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
		return nil
	case err := <-e:
		log.Fatalf("Error running server: %v", err)
	}

	return nil
}

const MaxContentLength = 2 * 1024 * 1024

// decode is a custom render.Decode function that makes sure the request doesn't exceed 2 megabytes.
func decode(r *http.Request, v any) error {
	if r.ContentLength > MaxContentLength {
		return server.APIError{
			Code: server.ErrRequestTooBig,
		}
	}

	// this is copied from render.Decode to replace r.Body with an io.LimitedReader
	var err error
	lr := io.LimitReader(r.Body, MaxContentLength)

	switch render.GetRequestContentType(r) {
	case render.ContentTypeJSON:
		err = render.DecodeJSON(lr, v)
	default:
		err = server.APIError{
			Code: server.ErrBadRequest,
		}
	}

	return err
}
