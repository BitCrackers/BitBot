package commands

import (
	"github.com/bwmarrin/discordgo"
)

type CommandKick struct{}

func (c *CommandKick) Name() string {
	return "kick"
}

func (c *CommandKick) Description() string {
	return "Kicks a user from the server."
}

func (c *CommandKick) AdminRequired() bool {
	return false
}

func (c *CommandKick) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
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
			Required:    true,
		},
	}
}

func (c *CommandKick) Exec(s *discordgo.Session, i *discordgo.InteractionCreate) {

	permissions, _ := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)

	if permissions&discordgo.PermissionKickMembers > 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "Kick called.",
			},
		})
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "Kick called but you don't have permissions in this channel to kick people.",
		},
	})
}
