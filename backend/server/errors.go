package server

import (
	"fmt"
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"github.com/go-chi/render"
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
	Details string `json:"details,omitempty"`

	RatelimitReset *int `json:"ratelimit_reset,omitempty"`

	// Status is set as the HTTP status code.
	Status int `json:"-"`
}

func (e APIError) Error() string {
	if e.Message == "" {
		e.Message = errCodeMessages[e.Code]
	}

	if e.Details != "" {
		return fmt.Sprintf("%s (code: %d) (%s)", e.Message, e.Code, e.Details)
	}

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
	ErrBadRequest          = 400
	ErrForbidden           = 403
	ErrNotFound            = 404
	ErrMethodNotAllowed    = 405
	ErrTooManyRequests     = 429
	ErrInternalServerError = 500 // catch-all code for unknown errors

	// Login/authorize error codes
	ErrInvalidState        = 1001
	ErrInvalidOAuthCode    = 1002
	ErrInvalidToken        = 1003 // a token was supplied, but it is invalid
	ErrInviteRequired      = 1004
	ErrInvalidTicket       = 1005 // invalid signup ticket
	ErrInvalidUsername     = 1006 // invalid username (when signing up)
	ErrUsernameTaken       = 1007 // username taken (when signing up)
	ErrInvitesDisabled     = 1008 // invites are disabled (unneeded)
	ErrInviteLimitReached  = 1009 // invite limit reached (when creating invites)
	ErrInviteAlreadyUsed   = 1010 // invite already used (when signing up)
	ErrDeletionPending     = 1011 // own user deletion pending, returned with undo code
	ErrRecentExport        = 1012 // latest export is too recent
	ErrUnsupportedInstance = 1013 // unsupported fediverse software
	ErrAlreadyLinked       = 1014 // user already has linked account of the same type

	// User-related error codes
	ErrUserNotFound = 2001

	// Member-related error codes
	ErrMemberNotFound     = 3001
	ErrMemberLimitReached = 3002
	ErrMemberNameInUse    = 3003
	ErrNotOwnMember       = 3004

	// General request error codes
	ErrRequestTooBig      = 4001
	ErrMissingPermissions = 4002
)

var errCodeMessages = map[int]string{
	ErrBadRequest:          "Bad request",
	ErrForbidden:           "Forbidden",
	ErrInternalServerError: "Internal server error",
	ErrNotFound:            "Not found",
	ErrTooManyRequests:     "Rate limit reached",
	ErrMethodNotAllowed:    "Method not allowed",

	ErrInvalidState:        "Invalid OAuth state",
	ErrInvalidOAuthCode:    "Invalid OAuth code",
	ErrInvalidToken:        "Supplied token was invalid",
	ErrInviteRequired:      "A valid invite code is required",
	ErrInvalidTicket:       "Invalid signup ticket",
	ErrInvalidUsername:     "Invalid username",
	ErrUsernameTaken:       "Username is already taken",
	ErrInvitesDisabled:     "Invites are disabled",
	ErrInviteLimitReached:  "Your account has reached the invite limit",
	ErrInviteAlreadyUsed:   "That invite code has already been used",
	ErrDeletionPending:     "Your account is pending deletion",
	ErrRecentExport:        "Your latest data export is less than 1 day old",
	ErrUnsupportedInstance: "Unsupported instance software",
	ErrAlreadyLinked:       "Your account is already linked to an account of this type",

	ErrUserNotFound: "User not found",

	ErrMemberNotFound:     "Member not found",
	ErrMemberLimitReached: "Member limit reached",
	ErrMemberNameInUse:    "Member name already in use",
	ErrNotOwnMember:       "Not your member",

	ErrRequestTooBig:      "Request too big (max 2 MB)",
	ErrMissingPermissions: "Your account or current token is missing required permissions for this action",
}

var errCodeStatuses = map[int]int{
	ErrBadRequest:          http.StatusBadRequest,
	ErrForbidden:           http.StatusForbidden,
	ErrInternalServerError: http.StatusInternalServerError,
	ErrNotFound:            http.StatusNotFound,
	ErrTooManyRequests:     http.StatusTooManyRequests,
	ErrMethodNotAllowed:    http.StatusMethodNotAllowed,

	ErrInvalidState:        http.StatusBadRequest,
	ErrInvalidOAuthCode:    http.StatusForbidden,
	ErrInvalidToken:        http.StatusUnauthorized,
	ErrInviteRequired:      http.StatusBadRequest,
	ErrInvalidTicket:       http.StatusBadRequest,
	ErrInvalidUsername:     http.StatusBadRequest,
	ErrUsernameTaken:       http.StatusBadRequest,
	ErrInvitesDisabled:     http.StatusForbidden,
	ErrInviteLimitReached:  http.StatusForbidden,
	ErrInviteAlreadyUsed:   http.StatusBadRequest,
	ErrDeletionPending:     http.StatusBadRequest,
	ErrRecentExport:        http.StatusBadRequest,
	ErrUnsupportedInstance: http.StatusBadRequest,
	ErrAlreadyLinked:       http.StatusBadRequest,

	ErrUserNotFound: http.StatusNotFound,

	ErrMemberNotFound:     http.StatusNotFound,
	ErrMemberLimitReached: http.StatusBadRequest,
	ErrMemberNameInUse:    http.StatusBadRequest,
	ErrNotOwnMember:       http.StatusForbidden,

	ErrRequestTooBig:      http.StatusBadRequest,
	ErrMissingPermissions: http.StatusForbidden,
}
