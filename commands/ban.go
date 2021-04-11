package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"time"
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
				Name:        "length",
				Description: "The length of the ban.",
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
	if permissions&discordgo.PermissionBanMembers < 0 {
		return
	}

	reason := fmt.Sprintf("Banned by: %s#%s.", i.Member.User.Username, i.Member.User.Discriminator)
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

	banTime := -1
	if timeString != "" {
		banTime, err = timeStringToSeconds(timeString)
		if err != nil {
			logrus.Errorf("%s: invalid time formatting", timeString)
			RespondWithError(s, i, fmt.Sprintf("%s: invalid time formatting", timeString))
			return
		}
	}

	err = ch.DB.BanUser(i.Data.Options[0].UserValue(s), i.Member.User, reason, banTime)
	if err != nil {
		logrus.Errorf("Error banning user: %v", err)
		RespondWithError(s, i, "Error adding user ban to database")
		return
	}

	err = s.GuildBanCreateWithReason(i.GuildID, i.Data.Options[0].UserValue(s).ID, reason, 0)

	if err != nil {
		logrus.Errorf("Error banning user: %v", err)
		RespondWithError(s, i, "Error banning user")
		return
	}

	banLength := "indefinite"
	if banTime != -1 {
		banLength = (time.Duration(banTime) * time.Second).String()
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf("**User %s#%s Banned**\n*Reason: %s*\n*Length: %s*", i.Data.Options[0].UserValue(s).Username, i.Data.Options[0].UserValue(s).Discriminator, reason, banLength),
		},
	})

	if err != nil {
		logrus.Errorf("Error responding to ban: %v", err)
	}
	return
}
