package commands

import (
	"fmt"
	"github.com/BitCrackers/BitBot/database"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (ch *CommandHandler) UnMuteCommand() *Command {
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

	if permissions&discordgo.PermissionKickMembers < 0 {
		return
	}

	u, err := ch.DB.GetUserRecord(i.Data.Options[0].UserValue(s))
	if err != nil {
		logrus.Errorf("Error fetching user record: %s\n", err)
		RespondWithError(s, i, "Couldn't fetch user record")
		return
	}
	if u.Mute.Empty() {
		RespondWithError(s, i, "User is not muted")
		return
	}

	u.Mute = database.Punishment{
		Type: -1,
	}
	err = ch.DB.SetUserRecord(u)
	if err != nil {
		logrus.Errorf("Error unmuting user: %s\n", err)
		RespondWithError(s, i, "Couldn't remove mute from database")
		return
	}

	err = s.GuildMemberRoleRemove(i.GuildID, i.Data.Options[0].UserValue(s).ID, ch.Config.MuteRoleId)
	if err != nil {
		logrus.Errorf("Error removing muted role from user: %s\n", err)
		RespondWithError(s, i, "Could not remove muted role from user")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf("**User %s#%s Unmuted", i.Data.Options[0].UserValue(s).Username, i.Data.Options[0].UserValue(s).Discriminator),
		},
	})
	if err != nil {
		logrus.Errorf("Error responding to unmute %v", err)
	}

	return
}
