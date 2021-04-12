package main

import (
	"github.com/BitCrackers/BitBot/commands"
	"github.com/BitCrackers/BitBot/config"
	"github.com/BitCrackers/BitBot/database"
	"github.com/BitCrackers/BitBot/modlog"
	"github.com/BitCrackers/BitBot/responses"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Check for required environment variables.
	token := os.Getenv("BITBOT_TOKEN")
	if token == "" {
		logrus.Fatalf("Couldn't read token from BITBOT_TOKEN environment variable")
	}

	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("Unable to load config: %v", err)
	}
	if cfg.GuildID == "" {
		panic("you were just saved from waiting 2 hours for discord to register slash commands globally")
	}

	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		logrus.Fatalf("Error while creating session: %v", err)
	}

	modlog, err := modlog.Create(&cfg, session)
	if err != nil {
		logrus.Fatalf("Unable to create modlog: %v", err)
	}

	db, err := database.New(session, &cfg, &modlog)
	if err != nil {
		logrus.Fatalf("Unable to start database: %v", err)
	}
	defer db.Close()

	cmdHandler := commands.CommandHandler{
		DB:     db,
		Config: &cfg,
		ModLog: &modlog,
	}

	// Setup session intents here. GuildMembers is needed for moderation slash commands.
	session.Identify.Intents = discordgo.IntentsAll

	// Create custom response handler for filters
	rh := responses.New()

	// Register handlers for filters.
	for _, filter := range cfg.Filters {
		handler, err := newFilterHandler(filter, &rh)
		if err != nil {
			logrus.Errorf("Error while creating filter: %v", err)
		}
		session.AddHandler(handler)
	}

	session.AddHandler(aumLog)

	if err = session.Open(); err != nil {
		logrus.Fatalf("Error while opening session: %v", err)
	}

	cmds := cmdHandler.Commands()

	// Contains all registered commands to avoid deleting commands where an error
	// occurred during registration and running into a panic.
	var registeredCmds []*commands.Command
	for _, cmd := range cmds {
		logrus.Infof("Registering slash command: %v", cmd.Name)
		if err = cmd.Register(session, cfg.GuildID); err != nil {
			logrus.Errorf("Error while registering %v command: %v", cmd.Name, err)
			continue
		}
		registeredCmds = append(registeredCmds, cmd)
	}

	// Wait until a termination signal is received.
	logrus.Infof("Bot is running successfully. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Clear slash commands.
	for _, cmd := range registeredCmds {
		logrus.Infof("Deleting slash command: %v", cmd.Name)
		if err = cmd.Delete(); err != nil {
			logrus.Errorf("Error while deleting command: %v", err)
		}
	}

	// Gracefully exit.
	session.Close()
}
