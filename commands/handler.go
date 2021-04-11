package commands

import (
	"fmt"
	"github.com/BitCrackers/BitBot/config"
	"github.com/BitCrackers/BitBot/database"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// CommandHandler provides a shared state for all command handler functions.
type CommandHandler struct {
	DB     *database.Database
	Config *config.Config
}

// Commands returns all commands of the CommandHandler. This should be updated
// to contain all `XxxxCommand()` methods of the handler.
func (ch *CommandHandler) Commands() []*Command {
	cmds := []*Command{
		ch.KickCommand(),
		ch.BanCommand(),
		ch.WarnCommand(),
		ch.BuildsCommand(),
		ch.MuteCommand(),
		ch.UnmuteCommand(),
	}
	if ch.Config.Debug {
		return append(
			cmds,
			ch.PingCommand(),
		)
	}
	return cmds
}

func RespondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, reason string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf("Command failed: %s", reason),
		},
	})
	if err != nil {
		logrus.Errorf("Error reporting error %v", err)
	}
}