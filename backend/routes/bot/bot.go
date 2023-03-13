package bot

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Bot struct {
	*server.Server

	publicKey ed25519.PublicKey
	baseURL   string
}

func (bot *Bot) UserAvatarURL(u db.User) string {
	if u.Avatar == nil {
		return ""
	}

	return bot.baseURL + "/media/users/" + u.ID.String() + "/" + *u.Avatar + ".webp"
}

func Mount(srv *server.Server, r chi.Router) {
	publicKey, err := hex.DecodeString(os.Getenv("DISCORD_PUBLIC_KEY"))
	if err != nil {
		return
	}

	b := &Bot{
		Server:    srv,
		publicKey: publicKey,
		baseURL:   os.Getenv("BASE_URL"),
	}

	r.HandleFunc("/interactions", b.handle)
}

func (bot *Bot) handle(w http.ResponseWriter, r *http.Request) {
	if !discordgo.VerifyInteraction(r, bot.publicKey) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var ev *discordgo.InteractionCreate

	if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	// we can always respond to ping with pong
	if ev.Type == discordgo.InteractionPing {
		log.Debug("received ping interaction")
		render.JSON(w, r, discordgo.InteractionResponse{
			Type: discordgo.InteractionResponsePong,
		})
		return
	}

	if ev.Type != discordgo.InteractionApplicationCommand {
		return
	}
	data := ev.ApplicationCommandData()

	switch data.Name {
	case "Show user's pronouns":
		bot.userPronouns(w, r, ev)
	case "Show author's pronouns":
	}
}

func (bot *Bot) userPronouns(w http.ResponseWriter, r *http.Request, ev *discordgo.InteractionCreate) {
	ctx := r.Context()

	var du *discordgo.User
	for _, user := range ev.ApplicationCommandData().Resolved.Users {
		du = user
		break
	}
	if du == nil {
		return
	}

	u, err := bot.DB.DiscordUser(ctx, du.ID)
	if err != nil {
		if err == db.ErrUserNotFound {
			respond(w, r, &discordgo.MessageEmbed{
				Description: du.String() + " does not have any pronouns set.",
			})
			return
		}

		log.Errorf("getting discord user: %v", err)
		return
	}

	avatarURL := du.AvatarURL("")
	if url := bot.UserAvatarURL(u); url != "" {
		avatarURL = url
	}
	name := u.Username
	if u.DisplayName != nil {
		name = fmt.Sprintf("%s (%s)", *u.DisplayName, u.Username)
	}
	url := bot.baseURL
	if url != "" {
		url += "/@" + u.Username
	}

	e := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    name,
			IconURL: avatarURL,
			URL:     url,
		},
	}

	if u.Bio != nil {
		e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
			Name:  "Bio",
			Value: *u.Bio,
		})
	}

	fields, err := bot.DB.UserFields(ctx, u.ID)
	if err != nil {
		respond(w, r, e)

		log.Errorf("getting user fields: %v", err)
		return
	}

	for _, field := range fields {
		var favs []db.FieldEntry

		for _, e := range field.Entries {
			if e.Status == db.StatusFavourite {
				favs = append(favs, e)
			}
		}

		if len(favs) == 0 {
			continue
		}

		var value string
		for _, fav := range favs {
			if len(fav.Value) > 500 {
				break
			}

			value += fav.Value + "\n"
		}

		e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
			Name:   field.Name,
			Value:  value,
			Inline: true,
		})
	}

	respond(w, r, e)
}

func respond(w http.ResponseWriter, r *http.Request, embeds ...*discordgo.MessageEmbed) {
	render.JSON(w, r, discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
			Flags:  uint64(discordgo.MessageFlagsEphemeral),
		},
	})
}
