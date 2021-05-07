package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (ch *CommandHandler) MuteCommand() *Command {
	return &Command{
		Name:        "mute",
		Description: "Mutes a user.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "user",
				Description: "The user to be muted.",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
			{
				Name:        "reason",
				Description: "The reason for kicking the user.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
			{
				Name:        "duration",
				Description: "The amount of time the mute should last.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
		HandlerFunc: ch.handleMute,
	}
}

func (ch *CommandHandler) handleMute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	permissions, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		logrus.Errorf("Error getting user permissions %v", err)
		RespondWithError(s, i, "Error fetching user permissions")
		return
	}

	for _, moderator := range ch.Config.Moderators {
		if i.Data.Options[0].UserValue(s).ID == moderator {
			logrus.Error("Cannot mute another moderator")
			RespondWithError(s, i, "Cannot mute another moderator")
			return
		}
	}

	if permissions&discordgo.PermissionKickMembers <= 0 || !ch.userIsModerator(i.Member.User.ID) {
		return
	}

	args := parseInteractionOptions(i.Data.Options)

	if ch.userIsModerator(args["user"].UserValue(s).ID) {
		RespondWithError(s, i, "cannot mute a moderator")
		return
	}

	reason := "unknown"
	if args["reason"] != nil && args["reason"].StringValue() != "" {
		reason = args["reason"].StringValue()
	}

	durationString := ""
	if args["duration"] != nil && args["duration"].StringValue() != "" {
		durationString = args["duration"].StringValue()
	}

	duration := time.Duration(-1)
	if durationString != "" {
		duration, err = timeStringToDuration(durationString)
		if err != nil {
			logrus.Errorf("%s: invalid time formatting", durationString)
			RespondWithError(s, i, fmt.Sprintf("%s: invalid time formatting", durationString))
			return
		}
	}

	if err = ch.DB.MuteUser(args["user"].UserValue(s).ID, i.Member.User.ID, reason, duration); err != nil {
		logrus.Errorf("Error while muting user: %v", err)
		RespondWithError(s, i, fmt.Sprintf("Error while muting user: %v", err))
		return
	}

	durationFmt := "indefinite"
	if duration != -1 {
		durationFmt = duration.String()
	}

	user := i.Data.Options[0].UserValue(s)
	err = ch.ModLog.SendEmbed(s, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("[MUTE] %s#%s", user.Username, user.Discriminator),
			IconURL: user.AvatarURL("256"),
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Color:     16753197,
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
				Value:  duration.String(),
				Inline: true,
			},
		},
	})
	if err != nil {
		logrus.Errorf("Error logging mute: %v", err)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf(
				"**User %s#%s Muted**\n*Reason: %s*\n*Length: %s*",
				i.Data.Options[0].UserValue(s).Username,
				i.Data.Options[0].UserValue(s).Discriminator,
				reason,
				durationFmt,
			),
		},
	})
	if err != nil {
		logrus.Errorf("Error responding to mute %v", err)
	}
}

func timeStringToDuration(t string) (time.Duration, error) {
	timeID := t[len(t)-1:]
	var multi time.Duration
	switch timeID {
	case "s":
		multi = time.Second
	case "m":
		multi = time.Minute
	case "h":
		multi = time.Hour
	case "d":
		multi = time.Hour * 24
	case "w":
		multi = time.Hour * 24 * 7
	case "y":
		multi = time.Hour * 24 * 365
	default:
		return -1, errors.New("invalid time multiplier")
	}
	timeString := strings.TrimRight(t, timeID)
	duration, err := strconv.Atoi(timeString)
	if err != nil {
		return -1, err
	}

	return time.Duration(duration) * multi, nil
}
