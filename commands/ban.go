package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (ch *CommandHandler) BanCommand() *Command {
	return &Command{
		Name:        "ban",
		Description: "Bans a user from the server.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "user",
				Description: "The user to be banned.",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
			{
				Name:        "reason",
				Description: "The reason for banning the user.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
			{
				Name:        "duration",
				Description: "The duration of the ban.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
		HandlerFunc: ch.handleBan,
	}
}

func (ch *CommandHandler) handleBan(s *discordgo.Session, i *discordgo.InteractionCreate) {
	permissions, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		logrus.Errorf("Error getting user permissions %v", err)
		RespondWithError(s, i, "Error fetching user permissions")
		return
	}

	if permissions&discordgo.PermissionBanMembers <= 0 {
		return
	}

	args := parseInteractionOptions(i.Data.Options)
	reason := fmt.Sprintf("Banned by: %s#%s.", i.Member.User.Username, i.Member.User.Discriminator)

	if args["reason"] != nil && args["reason"].StringValue() != "" {
		reason = args["reason"].StringValue()
	}

	if ch.userIsModerator(args["user"].UserValue(s).ID) {
		RespondWithError(s, i, "cannot ban a moderator")
		return
	}

	banTime := time.Duration(-1)
	if args["duration"] != nil && args["duration"].StringValue() != "" {
		banTime, err = timeStringToDuration(args["duration"].StringValue())
		if err != nil {
			logrus.Errorf("%s: invalid time formatting", args["duration"].StringValue())
			RespondWithError(s, i, fmt.Sprintf("%s: invalid time formatting", args["duration"].StringValue()))
			return
		}
	}

	if err = ch.DB.BanUser(args["user"].UserValue(s).ID, i.Member.User.ID, reason, banTime); err != nil {
		logrus.Errorf("Error banning user: %v", err)
		RespondWithError(s, i, "Error adding user ban to database")
		return
	}

	banLength := "indefinite"
	if banTime != -1 {
		banLength = banTime.String()
	}

	user := args["user"].UserValue(s)
	err = ch.ModLog.SendEmbed(s, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("[BAN] %s#%s", user.Username, user.Discriminator),
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
			{
				Name:   "Duration",
				Value:  banLength,
				Inline: true,
			},
		},
	})
	if err != nil {
		logrus.Errorf("Error logging ban: %v", err)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf("**User %s#%s Banned**\n*Reason: %s*\n*Length: %s*", args["user"].UserValue(s).Username, args["user"].UserValue(s).Discriminator, reason, banLength),
		},
	})

	if err != nil {
		logrus.Errorf("Error responding to ban: %v", err)
	}
	return
}
