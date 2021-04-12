package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (ch *CommandHandler) UnmuteCommand() *Command {
	return &Command{
		Name:        "unmute",
		Description: "Removes mute from a user.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "user",
				Description: "The user to be unmuted.",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
		},
		HandlerFunc: ch.handleUnMute,
	}
}

func (ch *CommandHandler) handleUnMute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	permissions, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		logrus.Errorf("Error getting user permissions %v", err)
		RespondWithError(s, i, "Error fetching user permissions")
		return
	}

	if permissions&discordgo.PermissionKickMembers <= 0 {
		return
	}

	u, err := ch.DB.UserRecord(i.Data.Options[0].UserValue(s).ID)
	if err != nil {
		logrus.Errorf("Error fetching user record: %s", err)
		RespondWithError(s, i, "Couldn't fetch user record")
		return
	}

	if u.Mute.Empty() {
		RespondWithError(s, i, "User is not muted")
		return
	}

	if err = ch.DB.UnmuteRecord(u, false); err != nil {
		logrus.Errorf("Error unmuting user: %s\n", err)
		RespondWithError(s, i, "Couldn't remove mute from database")
		return
	}

	user := i.Data.Options[0].UserValue(s)
	err = ch.ModLog.SendEmbed(s, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("[UNMUTE] %s#%s", user.Username, user.Discriminator),
			IconURL: user.AvatarURL("256"),
		},
		Description: "**Unmuted by moderator**",
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       3574686,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User",
				Value:  fmt.Sprintf("<@%s>", user.ID),
				Inline: true,
			},
			{
				Name:   "Moderator",
				Value:  fmt.Sprintf("<@%s>", i.Member.User.ID),
				Inline: true,
			},
		},
	})
	if err != nil {
		logrus.Errorf("Error logging unmute: %v", err)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf("**User %s#%s Unmuted**", i.Data.Options[0].UserValue(s).Username, i.Data.Options[0].UserValue(s).Discriminator),
		},
	})
	if err != nil {
		logrus.Errorf("Error responding to unmute %v", err)
	}

	return
}
