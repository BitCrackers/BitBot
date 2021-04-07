package commands

import (
	"fmt"
	"github.com/BitCrackers/BitBot/internal/router"
	"github.com/bwmarrin/discordgo"
)

var CommandKick = router.Command{
	Name:        "kick",
	Description: "Kicks a user from the server.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "user",
			Description: "The user to be kicked.",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
		{
			Name:        "reason",
			Description: "The reason for kicking the user.",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
	AdminRequired: true,
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		permissions, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if err != nil {
			fmt.Printf("error getting user permissions %v", err)
		}

		if permissions&discordgo.PermissionKickMembers > 0 {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: "Kick called.",
				},
			})
			if err != nil {
				fmt.Printf("error responding to kick %v", err)
			}

			return
		}
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "Kick called but you don't have permissions in this channel to kick people.",
			},
		})

		if err != nil {
			fmt.Printf("error responding to kick %v", err)
		}
	},
}
