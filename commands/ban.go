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
	if permissions&discordgo.PermissionBanMembers < 0 {
		return
	}

	args := parseInteractionOptions(i.Data.Options)
	reason := fmt.Sprintf("Banned by: %s#%s.", i.Member.User.Username, i.Member.User.Discriminator)

	if args["reason"] != nil && args["reason"].StringValue() != "" {
		reason = args["reason"].StringValue()
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

	err = ch.DB.BanUser(i.Data.Options[0].UserValue(s).ID, i.Member.User.ID, reason, banTime)
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
		banLength = banTime.String()
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
