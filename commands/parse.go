package commands

import (
	"fmt"
	"github.com/BitCrackers/BitBot/internal/router"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var CommandParse = router.Command{
	Name:        "parse",
	Description: "Parses a message. Debug command.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "sentence",
			Description: "The sentence that has to be parsed",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
	AdminRequired: false,
	Exec: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		sentence := i.Data.Options[0].StringValue()

		response := ""
		for _, a := range strings.Fields(sentence) {
			response += "["
			response += a
			response += "] "
		}

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{

			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: response,
			},
		})

		if err != nil {
			fmt.Printf("error trying to respond to parse command %v", err)
		}
	},
}
