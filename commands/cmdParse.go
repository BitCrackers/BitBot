package commands

import (
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (ch *CommandHandler) ParseCommand() *Command {
	return &Command{
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
		HandlerFunc: ch.handleParse,
	}
}

func (ch *CommandHandler) handleParse(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		logrus.Errorf("error trying to respond to parse command %v", err)
	}
}
