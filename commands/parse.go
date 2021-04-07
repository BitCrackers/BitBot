package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CommandParse struct{}

func (c *CommandParse) Name() string {
	return "parse"
}

func (c *CommandParse) Description() string {
	return "Parses a mesage. Debug command."
}

func (c *CommandParse) AdminRequired() bool {
	return false
}

func (c *CommandParse) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "sentence",
			Description: "The sentence that has to be parsed",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	}
}

func (c *CommandParse) Exec(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sentence := i.Data.Options[0].StringValue()

	response := ""
	for _, a := range strings.Fields(sentence) {
		response += "["
		response += a
		response += "] "
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{

		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: response,
		},
	})
}
