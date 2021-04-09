package main

import (
	"fmt"
	"github.com/BitCrackers/BitBot/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BitCrackers/BitBot/commands"
	"github.com/BitCrackers/BitBot/events"

	"github.com/BitCrackers/BitBot/internal/router"

	"github.com/bwmarrin/discordgo"
)

func main() {

	// Check for required environment variables.
	err := config.Load()
	bbToken := os.Getenv("BITBOT_TOKEN")

	if err != nil {
		log.Fatalf("Unable to load config: %v", err)
	}

	// Just for fun at the moment, but we should probably only do this if $BITBOT_DEBUG is true.
	// TODO: Set up with debug env. variable.
	fmt.Println("$BITBOT_TOKEN: ", bbToken)
	fmt.Println("$BITBOT_GUILDID: ", config.C.GuildID)

	bot, err := discordgo.New("Bot " + bbToken)

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

	if config.C.Debug {
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
	cmdHandler.CreateCommands(bot, config.C.GuildID)

	// Wait until a termination signal is received.
	fmt.Println("Bot is running successfully. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	//Clear slash commands.
	cmdHandler.ClearCommands(bot, config.C.GuildID)

	// Gracefully exit.
	_ = bot.Close()
}
