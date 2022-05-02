package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"gitlab.com/1f320/pronouns/backend/log"
)

// WrapHandler wraps a modified http.HandlerFunc into a stdlib-compatible one.
// The inner HandlerFunc additionally returns an error.
func WrapHandler(hn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := hn(w, r)
		if err != nil {
			// if the function returned an API error, just render that verbatim
			// we can assume that it also logged the error (if that was needed)
			if apiErr, ok := err.(APIError); ok {
				apiErr.prepare()

				render.Status(r, apiErr.Status)
				render.JSON(w, r, apiErr)
				return
			}

			// otherwise, we log the error and return an internal server error message
			log.Errorf("error in http handler: %v", err)

			apiErr := APIError{Code: ErrInternalServerError}
			apiErr.prepare()

			render.Status(r, apiErr.Status)
			render.JSON(w, r, apiErr)
		}
	}
}

// APIError is an object returned by the API when an error occurs.
// It implements the error interface and can be returned by handlers.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`

	// Status is
	Status int `json:"-"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s (code: %d)", e.Message, e.Code)
}

func (e *APIError) prepare() {
	if e.Status == 0 {
		e.Status = errCodeStatuses[e.Code]
	}

	if e.Message == "" {
		e.Message = errCodeMessages[e.Code]
	}
}

// Error code constants
const (
	ErrInternalServerError = 500 // catch-all code for unknown errors
)

var errCodeMessages = map[int]string{
	ErrInternalServerError: "Internal server error",
}

var errCodeStatuses = map[int]int{
	ErrInternalServerError: http.StatusInternalServerError,
}
