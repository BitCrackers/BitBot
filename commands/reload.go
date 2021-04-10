package commands

import (
	"fmt"

	"github.com/BitCrackers/BitBot/internal/config"
	"github.com/BitCrackers/BitBot/internal/router"
	"github.com/bwmarrin/discordgo"
)

var CommandReload = router.Command{
	Name:        "reload",
	Description: "Reloads the config file.",
	Options:     make([]*discordgo.ApplicationCommandOption, 0),
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		for _, m := range config.C.Moderators {
			fmt.Printf("Checking reload perm: %s - %s\n", m, i.Member.User.ID)
			if m == i.Member.User.ID {
				config.Load()
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Content: "Config file reloaded.",
					},
				})

				if err != nil {
					fmt.Printf("error trying to respond to ping command %v", err)
				}

				break
			}
		}
	},
}
