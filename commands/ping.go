package commands

import (
	"fmt"
	"github.com/BitCrackers/BitBot/internal/router"
	"github.com/bwmarrin/discordgo"
)

var CommandPing = router.Command{
	Name:          "ping",
	Description:   "Pong!",
	Options:       make([]*discordgo.ApplicationCommandOption, 0),
	AdminRequired: false,
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "Pong",
			},
		})

		if err != nil {
			fmt.Printf("error trying to respond to ping command %v", err)
		}
	},
}
