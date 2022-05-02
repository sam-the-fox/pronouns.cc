package auth

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"emperror.dev/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/xid"
	"gitlab.com/1f320/pronouns/backend/log"
)

// Claims are the claims used in a token.
type Claims struct {
	UserID xid.ID `json:"sub"`

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

const expireDays = 30

// CreateToken creates a token for the given user ID.
// It expires after 30 days.
func (v *Verifier) CreateToken(userID xid.ID) (string, error) {
	now := time.Now()
	expires := now.Add(expireDays * 24 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "pronouns",
			ExpiresAt: jwt.NewNumericDate(expires),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	})

	return token.SignedString(v.key)
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
