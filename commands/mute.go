package commands

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"math"
	"strconv"
	"strings"
	"time"
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
				Name:        "length",
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

	if permissions&discordgo.PermissionKickMembers < 0 {
		return
	}

	reason := "unknown"
	timeString := ""
	if len(i.Data.Options) == 2 {
		if i.Data.Options[1].Name == "reason" {
			reason = i.Data.Options[1].StringValue()
		} else {
			timeString = i.Data.Options[1].StringValue()
		}
	}
	if len(i.Data.Options) == 3 {
		reason = i.Data.Options[1].StringValue()
		timeString = i.Data.Options[2].StringValue()
	}

	muteTime := -1
	if timeString != "" {
		muteTime, err = timeStringToSeconds(timeString)
		if err != nil {
			logrus.Errorf("%s: invalid time formatting", timeString)
			RespondWithError(s, i, fmt.Sprintf("%s: invalid time formatting", timeString))
			return
		}
	}
	u, err := ch.DB.GetUserRecord(i.Data.Options[0].UserValue(s))
	if err != nil {
		logrus.Errorf("Error fetching user record: %s\n", err)
		RespondWithError(s, i, "Could fetch user record")
		return
	}

	if !u.Mute.Empty() {
		muteExpire := "never"
		if u.Mute.Length != -1 {
			muteExpire = u.Mute.Date.Add(time.Duration(u.Mute.Length)).Sub(u.Mute.Date).String()
		}
		logrus.Errorf("User is already muted\n*Mute expires:%s*", muteExpire)
		RespondWithError(s, i, fmt.Sprintf("User is already muted\n*Mute expires: %s*", muteExpire))
		return
	}

	err = ch.DB.MuteUser(i.Data.Options[0].UserValue(s), i.Member.User, reason, muteTime)
	if err != nil {
		logrus.Errorf("Error warning user: %s\n", err)
		RespondWithError(s, i, "Could not add muted user to database")
		return
	}

	err = s.GuildMemberRoleAdd(i.GuildID, i.Data.Options[0].UserValue(s).ID, ch.Config.MuteRoleId)
	if err != nil {
		logrus.Errorf("Error giving user muted role: %s\n", err)
		RespondWithError(s, i, "Could not add muted role to user")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf("**User %s#%s Muted**\n*Reason: %s*\n*Length: %s*", i.Data.Options[0].UserValue(s).Username, i.Data.Options[0].UserValue(s).Discriminator, reason, secondsToString(muteTime)),
		},
	})
	if err != nil {
		logrus.Errorf("Error responding to mute %v", err)
	}

	return
}

func timeStringToSeconds(t string) (int, error) {
	timeID := t[len(t)-1:]
	multi := 1
	switch timeID {
	case "s":
		break
	case "m":
		multi = 60
		break
	case "h":
		multi = 3600
		break
	case "d":
		multi = 86400
		break
	case "w":
		multi = 604800
		break
	case "y":
		multi = 31622400
		break
	default:
		return -1, errors.New("invalid time multiplier")
	}
	timeString := strings.TrimRight(t, timeID)
	time, err := strconv.Atoi(timeString)
	if err != nil {
		return -1, err
	}

	return time * multi, nil
}

//https://www.socketloop.com/tutorials/golang-convert-seconds-to-human-readable-time-format-example
func plural(count int, singular string) (result string) {
	if count == 1 {
		result = strconv.Itoa(count) + " " + singular + " "
	} else {
		result = strconv.Itoa(count) + " " + singular + "s "
	}
	return
}

func secondsToString(input int) (result string) {
	if input == -1 {
		result = "indefinite"
		return
	}
	years := math.Floor(float64(input) / 60 / 60 / 24 / 7 / 30 / 12)
	seconds := input % (60 * 60 * 24 * 7 * 30 * 12)
	months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 30)
	seconds = input % (60 * 60 * 24 * 7 * 30)
	weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
	seconds = input % (60 * 60 * 24 * 7)
	days := math.Floor(float64(seconds) / 60 / 60 / 24)
	seconds = input % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = input % 60

	if years > 0 {
		result = plural(int(years), "year") + plural(int(months), "month") + plural(int(weeks), "week") + plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(seconds, "second")
	} else if months > 0 {
		result = plural(int(months), "month") + plural(int(weeks), "week") + plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(seconds, "second")
	} else if weeks > 0 {
		result = plural(int(weeks), "week") + plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(seconds, "second")
	} else if days > 0 {
		result = plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else if hours > 0 {
		result = plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(seconds, "second")
	} else if minutes > 0 {
		result = plural(int(minutes), "minute") + plural(seconds, "second")
	} else {
		result = plural(seconds, "second")
	}
	result = strings.TrimSpace(result)

	return
}
