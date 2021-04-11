package commands

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

func (ch *CommandHandler) KickCommand() *Command {
	return &Command{
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
				Required:    false,
			},
		},
		HandlerFunc: ch.handleKick,
	}
}

func (ch *CommandHandler) handleKick(s *discordgo.Session, i *discordgo.InteractionCreate) {
	permissions, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		logrus.Errorf("Error getting user permissions %v", err)
		RespondWithError(s, i, "Error fetching user permissions")
		return
	}

	if permissions&discordgo.PermissionKickMembers < 0 {
		return
	}

	var reason string
	if len(i.Data.Options) > 1 {
		reason = i.Data.Options[1].StringValue()
	} else {
		reason = fmt.Sprintf("Kicked by: %s#%s.", i.Member.User.Username, i.Member.User.Discriminator)
	}

	err = s.GuildMemberDeleteWithReason(i.GuildID, i.Data.Options[0].UserValue(s).ID, reason)

	if err != nil {
		logrus.Errorf("Error kicking user: %v", err)
		RespondWithError(s, i, "Error kicking user")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf("**User %s#%s Kicked**\n*Reason: %s*", i.Data.Options[0].UserValue(s).Username, i.Data.Options[0].UserValue(s).Discriminator, reason),
		},
	})
	if err != nil {
		logrus.Errorf("Error responding to kick %v", err)
	}

	return
}
