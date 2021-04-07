package commands

import (
	"github.com/bwmarrin/discordgo"
)

type CommandBan struct{}

func (c *CommandBan) Name() string {
	return "ban"
}

func (c *CommandBan) Description() string {
	return "Bans a user from the server."
}

func (c *CommandBan) AdminRequired() bool {
	return false
}

func (c *CommandBan) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
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
			Required:    true,
		},
	}
}

func (c *CommandBan) Exec(s *discordgo.Session, i *discordgo.InteractionCreate) {

	permissions, _ := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)

	if permissions&discordgo.PermissionBanMembers > 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "Ban called.",
			},
		})
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "Ban called but you don't have permissions in this channel to ban people.",
		},
	})
}
