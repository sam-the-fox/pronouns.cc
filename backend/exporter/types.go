package exporter

import (
	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"github.com/rs/xid"
)

type userExport struct {
	ID          xid.ID  `json:"id"`
	Username    string  `json:"name"`
	DisplayName *string `json:"display_name"`
	Bio         *string `json:"bio"`

	Links []string `json:"links"`

	Names    []db.FieldEntry   `json:"names"`
	Pronouns []db.PronounEntry `json:"pronouns"`
	Fields   []db.Field        `json:"fields"`

	Fediverse         *string `json:"fediverse"`
	FediverseUsername *string `json:"fediverse_username"`
	FediverseInstance *string `json:"fediverse_instance"`

	Discord         *string `json:"discord"`
	DiscordUsername *string `json:"discord_username"`

	Tumblr         *string `json:"tumblr"`
	TumblrUsername *string `json:"tumblr_username"`

	Google         *string `json:"google"`
	GoogleUsername *string `json:"google_username"`

	MaxInvites int `json:"max_invites"`

	Warnings []db.Warning `json:"warnings"`
}

func dbUserToExport(u db.User, fields []db.Field, warnings []db.Warning) userExport {
	return userExport{
		ID:                u.ID,
		Username:          u.Username,
		DisplayName:       u.DisplayName,
		Bio:               u.Bio,
		Links:             db.NotNull(u.Links),
		Names:             db.NotNull(u.Names),
		Pronouns:          db.NotNull(u.Pronouns),
		Fields:            db.NotNull(fields),
		Discord:           u.Discord,
		DiscordUsername:   u.DiscordUsername,
		Tumblr:            u.Tumblr,
		TumblrUsername:    u.TumblrUsername,
		Google:            u.Google,
		GoogleUsername:    u.GoogleUsername,
		MaxInvites:        u.MaxInvites,
		Fediverse:         u.Fediverse,
		FediverseUsername: u.FediverseUsername,
		FediverseInstance: u.FediverseInstance,
		Warnings:          db.NotNull(warnings),
	}
}

type memberExport struct {
	ID          xid.ID            `json:"id"`
	Name        string            `json:"name"`
	DisplayName *string           `json:"display_name"`
	Bio         *string           `json:"bio"`
	Links       []string          `json:"links"`
	Names       []db.FieldEntry   `json:"names"`
	Pronouns    []db.PronounEntry `json:"pronouns"`
	Fields      []db.Field        `json:"fields"`
	Unlisted    bool              `json:"unlisted"`
}

func dbMemberToExport(m db.Member, fields []db.Field) memberExport {
	return memberExport{
		ID:          m.ID,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Bio:         m.Bio,
		Links:       db.NotNull(m.Links),
		Names:       db.NotNull(m.Names),
		Pronouns:    db.NotNull(m.Pronouns),
		Fields:      db.NotNull(fields),
		Unlisted:    m.Unlisted,
	}
}
