package commands

import (
	"fmt"
	"github.com/BitCrackers/BitBot/database"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (ch *CommandHandler) ReactionRoleCommand() *Command {
	return &Command{
		Name:        "reactionrole",
		Description: "Creates a message that can be reacted to to receive roles.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "channel",
				Description: "Where the message will be posted",
				Type:        discordgo.ApplicationCommandOptionChannel,
				Required:    true,
			},
			{
				Name:        "embed",
				Description: "Whether the message will be an embed or not.",
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Required:    true,
			},
			{
				Name:        "message",
				Description: "The message that will be sent.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "emote",
				Description: "The reaction emote that will be sent along with the message.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "role",
				Description: "The role that will be given when reacted to the message.",
				Type:        discordgo.ApplicationCommandOptionRole,
				Required:    true,
			},
		},
		HandlerFunc: ch.handleReactionCommand,
	}
}

func (ch *CommandHandler) handleReactionCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	permissions, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		logrus.Errorf("Error getting user permissions %v", err)
		RespondWithError(s, i, "Error fetching user permissions")
	}

	if permissions&discordgo.PermissionKickMembers <= 0 || !ch.userIsModerator(i.Member.User.ID) {
		return
	}

	args := parseInteractionOptions(i.Data.Options)

	g, err := s.Guild(ch.Config.GuildID)
	if err != nil {
		logrus.Errorf("Error fetching guild by id %v", err)
		RespondWithError(s, i, "Error fetching guild by id")
	}
	var m *discordgo.Message

	if args["embed"].BoolValue() {
		e := discordgo.MessageEmbed{
			Description: args["message"].StringValue(),
		}
		m, err = s.ChannelMessageSendEmbed(args["channel"].ChannelValue(s).ID, &e)
	} else {
		m, err = s.ChannelMessageSend(args["channel"].ChannelValue(s).ID, args["message"].StringValue())
	}

	if err != nil {
		logrus.Errorf("Error sending role reaction message %v", err)
		RespondWithError(s, i, "Error sending role reaction message")
	}

	r := database.ReactionRole{
		ID:    m.ID,
		Emote: args["emote"].StringValue(),
		Role:  args["role"].RoleValue(s, g.ID).ID,
		Channel: args["channel"].ChannelValue(s).ID,
	}
	err = ch.DB.AddReactionRole(r.ID, r)
	if err != nil {
		logrus.Errorf("Error adding reaction role to database %v", err)
		RespondWithError(s, i, "Error adding reaction role to database")
	}

	err = s.MessageReactionAdd(args["channel"].ChannelValue(s).ID, m.ID, args["emote"].StringValue())
	if err != nil {
		logrus.Errorf("Error adding reaction to message %v", err)
		RespondWithError(s, i, "Error adding reaction to message")
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf("Message sent"),
		},
	})
	if err != nil {
		logrus.Errorf("Error responding to reactionrole %v", err)
	}

	return
}
