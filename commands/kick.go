package commands

import (
	"fmt"
	"time"

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

	if permissions&discordgo.PermissionKickMembers <= 0 || !ch.userIsModerator(i.Member.User.ID) {
		return
	}

	args := parseInteractionOptions(i.Data.Options)
	reason := fmt.Sprintf("kicked by: %s#%s.", i.Member.User.Username, i.Member.User.Discriminator)

	if args["reason"] != nil && args["reason"].StringValue() != "" {
		reason = args["reason"].StringValue()
	}

	if ch.userIsModerator(args["user"].UserValue(s).ID) {
		RespondWithError(s, i, "cannot kick a moderator")
		return
	}

	if err = s.GuildMemberDeleteWithReason(i.GuildID, args["user"].UserValue(s).ID, reason); err != nil {
		logrus.Errorf("Error kicking user: %v", err)
		RespondWithError(s, i, "Error kicking user")
		return
	}

	user := args["user"].UserValue(s)
	err = ch.ModLog.SendEmbed(s, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("[KICK] %s#%s", user.Username, user.Discriminator),
			IconURL: user.AvatarURL("256"),
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Color:     16754451,
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
			{
				Name:   "Reason",
				Value:  reason,
				Inline: true,
			},
		},
	})

	if err != nil {
		logrus.Errorf("Error logging kick: %v", err)
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
