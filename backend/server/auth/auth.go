package auth

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"emperror.dev/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/xid"
)

// Claims are the claims used in a token.
type Claims struct {
	UserID      xid.ID `json:"sub"`
	TokenID     xid.ID `json:"jti"`
	UserIsAdmin bool   `json:"adm"`

	// APIToken specifies whether this token was generated for the API or for the website.
	// API tokens cannot perform some destructive actions, such as DELETE /users/@me.
	APIToken bool `json:"atn"`
	// TokenWrite specifies whether this token can be used for write actions.
	// If set to false, this token can only be used for read actions.
	TokenWrite bool `json:"twr"`

	jwt.RegisteredClaims
}

type Verifier struct {
	key []byte
}

func New() *Verifier {
	raw := os.Getenv("HMAC_KEY")
	if raw == "" {
		log.Fatal("$HMAC_KEY is not set")
	}

	key, err := base64.URLEncoding.DecodeString(raw)
	if err != nil {
		log.Fatal("$HMAC_KEY is not a valid base 64 string")
	}

	return &Verifier{key: key}
}

// ExpireDays is after how many days the token will expire.
const ExpireDays = 30

// CreateToken creates a token for the given user ID.
// It expires after 30 days.
func (v *Verifier) CreateToken(userID, tokenID xid.ID, isAdmin bool, isAPIToken bool, isWriteToken bool) (token string, err error) {
	now := time.Now()
	expires := now.Add(ExpireDays * 24 * time.Hour)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID:      userID,
		TokenID:     tokenID,
		UserIsAdmin: isAdmin,
		APIToken:    isAPIToken,
		TokenWrite:  isWriteToken,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "pronouns",
			ExpiresAt: jwt.NewNumericDate(expires),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	})

	return t.SignedString(v.key)
}

// Claims parses the given token and returns its Claims.
// If the token is invalid, returns an error.
func (v *Verifier) Claims(token string) (c Claims, err error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(`unexpected signing method "%v"`, t.Header["alg"])
		}

		return v.key, nil
	})

	if err != nil {
		return c, errors.Wrap(err, "parsing token")
	}

	if c, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		return *c, nil
	}

	return c, fmt.Errorf("unknown claims type %T", parsed.Claims)
}
