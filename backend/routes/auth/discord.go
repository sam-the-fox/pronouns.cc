package auth

import (
	"os"

	"golang.org/x/oauth2"
)

var oauthConfig = oauth2.Config{
	ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
	ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
	Endpoint: oauth2.Endpoint{
		AuthURL:   "https://discord.com/api/oauth2/authorize",
		TokenURL:  "https://discord.com/api/oauth2/token",
		AuthStyle: oauth2.AuthStyleInParams,
	},
	Scopes: []string{"identify"},
}
