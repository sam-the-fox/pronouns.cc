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

	Discord         *string `json:"discord"`
	DiscordUsername *string `json:"discord_username"`

	MaxInvites int `json:"max_invites"`
}

func dbUserToExport(u db.User, fields []db.Field) userExport {
	return userExport{
		ID:              u.ID,
		Username:        u.Username,
		DisplayName:     u.DisplayName,
		Bio:             u.Bio,
		Links:           u.Links,
		Names:           u.Names,
		Pronouns:        u.Pronouns,
		Fields:          fields,
		Discord:         u.Discord,
		DiscordUsername: u.DiscordUsername,
		MaxInvites:      u.MaxInvites,
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
}

func dbMemberToExport(m db.Member, fields []db.Field) memberExport {
	return memberExport{
		ID:          m.ID,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Bio:         m.Bio,
		Links:       m.Links,
		Names:       m.Names,
		Pronouns:    m.Pronouns,
		Fields:      fields,
	}
}
