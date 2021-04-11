package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (ch *CommandHandler) WarnCommand() *Command {
	return &Command{
		Name:        "warn",
		Description: "Warns a user.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "user",
				Description: "The user to be warned.",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
			{
				Name:        "reason",
				Description: "The reason for kicking the user.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
		HandlerFunc: ch.handleWarn,
	}
}

func (ch *CommandHandler) handleWarn(s *discordgo.Session, i *discordgo.InteractionCreate) {
	permissions, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		logrus.Errorf("Error getting user permissions %v", err)
		RespondWithError(s, i, "Error fetching user permissions")
	}

	if permissions&discordgo.PermissionKickMembers < 0 {
		return
	}

	var reason string
	if len(i.Data.Options) > 1 {
		reason = i.Data.Options[1].StringValue()
	} else {
		reason = "unknown"
	}

	err = ch.DB.WarnUser(i.Data.Options[0].UserValue(s), i.Member.User, reason)

	if err != nil {
		fmt.Printf("Error warning user: %s\n", err)
		RespondWithError(s, i, "Error warning user")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf("**User %s#%s Warned**\n*Reason: %s*", i.Data.Options[0].UserValue(s).Username, i.Data.Options[0].UserValue(s).Discriminator, reason),
		},
	})
	if err != nil {
		logrus.Errorf("Error responding to warn %v", err)
	}

	return
}
