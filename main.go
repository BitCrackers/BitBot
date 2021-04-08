package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/BitCrackers/BitBot/commands"
	"github.com/BitCrackers/BitBot/events"

	"github.com/BitCrackers/BitBot/helpers"
	"github.com/BitCrackers/BitBot/internal/router"

	"github.com/bwmarrin/discordgo"
)

func main() {

	// Check for required environment variables.
	bbToken := os.Getenv("BITBOT_TOKEN")
	bbDebug := os.Getenv("BITBOT_DEBUG")
	bbOwner := os.Getenv("BITBOT_OWNERID")
	bbGuild := os.Getenv("BITBOT_GUILDID")

	d, _ := strconv.ParseBool(bbDebug)

	// Set environment settings.
	helpers.SetSettings(bbToken, d, bbOwner, bbGuild)

	// Throw specific exit codes along with a helpful message.
	if bbToken == "" {
		fmt.Println("> $BITBOT_TOKEN has not been exported.")
		os.Exit(2)
	}

	if bbDebug == "" {
		fmt.Println("> $BITBOT_DEBUG has not been exported.")
		os.Exit(3)
	}

	if bbOwner == "" {
		fmt.Println("> $BITBOT_OWNERID has not been exported.")
		os.Exit(4)
	}

	// TODO: If this check fails, default to normal commands.
	if bbGuild == "" {
		fmt.Println("> $BITBOT_GUILDID has not been exported.")
		os.Exit(5)
	}

	// Just for fun at the moment, but we should probably only do this if $BITBOT_DEBUG is true.
	// TODO: Set up with debug env. variable.
	fmt.Println("$BITBOT_TOKEN: ", bbToken)
	fmt.Println("$BITBOT_DEBUG: ", bbDebug)
	fmt.Println("$BITBOT_OWNERID: ", bbOwner)
	fmt.Println("$BITBOT_GUILDID: ", bbGuild)

	setupBot(bbToken, bbGuild)
}

func setupBot(token string, guildId string) {

	// Create a Discord session.
	bot, err := discordgo.New("Bot " + token)

	cmdHandler := router.NewCommandHandler()

	if err != nil {
		fmt.Println("> ", err)
		os.Exit(5)
	}

	// Add event handlers here.
	bot.AddHandler(events.NewMessageHandler().Handler)
	bot.AddHandler(events.NewEditHandler().Handler)
	bot.AddHandler(events.NewDeleteHandler().Handler)

	// Set up command handler.
	bot.AddHandler(cmdHandler.Handler) // Add commands here.

	if helpers.GetSettings().Debug {
		cmdHandler.RegisterCommand(commands.CommandPing)
		cmdHandler.RegisterCommand(commands.CommandParse)
	}

	cmdHandler.RegisterCommand(commands.CommandKick)
	cmdHandler.RegisterCommand(commands.CommandBan)

	// Setup bot intents here. GuildMembers is needed for moderation slash commands.
	// bot.Identify.Intents = discordgo.IntentsAll
	bot.Identify.Intents = discordgo.IntentsAll

	err = bot.Open()

	if err != nil {
		fmt.Println("> ", err)
		os.Exit(6)
	}

	//Create all the slash commands here as it can only be done after the bot starts
	cmdHandler.CreateCommands(bot, guildId)

	// Wait until a termination signal is received.
	fmt.Println("Bot is running successfully. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	//Clear slash commands.
	cmdHandler.ClearCommands(bot, guildId)

	// Gracefully exit.
	_ = bot.Close()
}
