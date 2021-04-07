package commands

import (
	"fmt"

	"github.com/BitCrackers/BitBot/internal/router"
	"github.com/bwmarrin/discordgo"
)

var CommandBan = router.Command{
	Name:        "ban",
	Description: "Bans a user from the server.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "user",
			Description: "The user to be banned.",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
		{
			Name:        "reason",
			Description: "The reason for banning the user.",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
		{
			Name:        "days",
			Description: "The number of days worth of user messages to delete.",
			Type:        discordgo.ApplicationCommandOptionInteger,
			Required:    true,
		},
	},
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		permissions, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if err != nil {
			fmt.Printf("Error getting user permissions %s", err.Error())
		}

		var reason string
		if len(i.Data.Options) > 1 {
			reason = i.Data.Options[1].StringValue()
		} else {
			reason = fmt.Sprintf("Banned by: %s#%s.", i.Member.User.Username, i.Member.User.Discriminator)
		}
		if permissions&discordgo.PermissionBanMembers > 0 {

			err := s.GuildBanCreateWithReason(i.GuildID, i.Data.Options[0].UserValue(s).ID, reason, int(i.Data.Options[2].IntValue()))

			if err != nil {
				fmt.Printf("Error banning user: %s", err.Error())
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: fmt.Sprintf("**User %s#%s Banned**\n*Reason: %s*", i.Data.Options[0].UserValue(s).Username, i.Data.Options[0].UserValue(s).Discriminator, reason),
				},
			})

			if err != nil {
				fmt.Printf("Error responding to ban %s", err.Error())
			}
			return
		}
	},
}
