package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BitCrackers/BitBot/events"

	"github.com/bwmarrin/discordgo"
)

func main() {

	// Check for required environment variables.
	bbToken := os.Getenv("BITBOT_TOKEN")
	bbDebug := os.Getenv("BITBOT_DEBUG")
	bbOwner := os.Getenv("BITBOT_OWNERID")

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

	// Just for fun at the moment, but we should probably only do this if $BITBOT_DEBUG is true.
	// TODO: Set up with debug env. variable.
	fmt.Println("$BITBOT_TOKEN: ", bbToken)
	fmt.Println("$BITBOT_DEBUG: ", bbDebug)
	fmt.Println("$BITBOT_OWNERID: ", bbOwner)

	setupBot(bbToken)
}

func setupBot(token string) {

	// Create a Discord session.
	bot, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("> ", err)
		os.Exit(5)
	}

	// Add event handlers here.
	bot.AddHandler(events.OnMessageCreate)

	// Setup bot intents here. For now I just have it as Messages, but we can switch to All easily enough.
	// bot.Identify.Intents = discordgo.IntentsAll
	bot.Identify.Intents = discordgo.IntentsGuildMessages

	err = bot.Open()

	if err != nil {
		fmt.Println("> ", err)
		os.Exit(6)
	}

	// Wait until a termination signal is received.
	fmt.Println("Bot is running successfully. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Gracefull exit.
	bot.Close()
}
