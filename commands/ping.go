package commands

import (
	"github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

func (ch *CommandHandler) PingCommand() *Command {
	return &Command{
		Name:        "ping",
		Description: "Pong!",
		Options:     []*discordgo.ApplicationCommandOption{},
		HandlerFunc: ch.handlePing,
	}
}

func (ch *CommandHandler) handlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "Pong",
		},
	})

	if err != nil {
		logrus.Errorf("error trying to respond to ping command %v", err)
	}
}