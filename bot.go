package sento

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// OpenCloser interface used mainly for callbacks
type OpenCloser interface {
	Open() error
	Close() error
}

// Bot is a sento-powered bot application
type Bot struct {
	sess      *discordgo.Session
	handlers  map[string]Handler
	listeners map[EventType][]EventListener
	cfg       *Config
	logger    *zap.Logger
}

// New returns a new sento-powered discord bot
func New(options ...Option) (bot *Bot, err error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	bot = &Bot{
		logger:    logger,
		listeners: make(map[EventType][]EventListener),
	}
	for _, op := range options {
		err = op(bot)
		if err != nil {
			break
		}
	}
	return
}

// Start an instance of the bot
func (bot *Bot) Start() (err error) {
	bot.LogInfo("Creating session")
	bot.sess, err = discordgo.New("Bot " + bot.cfg.Token)
	if err != nil {
		// TODO: Maybe modify error message
		// Could not connect to host/discord
		bot.LogInfo("Error creating session")
		return err
	}

	// Add handlers
	bot.LogInfo("Starting all handlers")
	for _, handler := range bot.handlers {
		err = handler.Start(bot)
		if err != nil {
			// TODO: Maybe modify error message
			// Error while starting handler
			return err
		}
	}
	bot.sess.AddHandler(bot.handleCreateMessage)

	bot.LogInfo("Opening the connection")
	err = bot.sess.Open()
	if err != nil {
		bot.LogInfo("Error opening the connection")
		return
	}

	bot.LogInfo("Listening . . .")
	bot.EmitEvent(EventConnected, nil)

	return
}

// Stop the bot
func (bot *Bot) Stop() (err error) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.LogInfo("Signal received")
	bot.LogInfo("Stopping all handlers")
	for _, handler := range bot.handlers {
		err = handler.Stop(bot)
		if err != nil {
			// TODO: Maybe modify error message
			// Error while stoping handler
			return err
		}
	}

	bot.LogInfo("Closing connection")
	err = bot.sess.Close()
	if err == nil {
		bot.LogInfo("Connection closed")
	}
	bot.EmitEvent(EventDisconnected, nil)
	return
}

func (bot *Bot) handleCreateMessage(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == sess.State.User.ID {
		return // Ignore messages sent by this bot
	}

	// TODO: Fetch prefix from database
	prefix := ""

	if prefix == "" {
		// If there is no prefix for the server
		// use the default
		prefix = DefaultConfig.Prefix
	}

	if !strings.HasPrefix(msg.Content, prefix) {
		// Ignore messages without prefix
		return
	}

	// Grab the trigger
	triggerEnd := strings.Index(msg.Content, " ")
	if triggerEnd == -1 {
		triggerEnd = len(msg.Content)
	}
	trigger := msg.Content[len(prefix):triggerEnd]

	handler, triggerExists := bot.handlers[trigger]
	if !triggerExists {
		// Ignore messages with no handlers
		return
	}

	// Handle message
	handleInfo := HandleInfo{
		AuthorID:  msg.Author.ID,
		ChannelID: msg.ChannelID,
		GuildID:   msg.GuildID,
		Trigger:   trigger,
	}
	err := handler.Handle(bot, handleInfo)

	// TODO: Maybe make it prettier?
	logFields := []LogField{
		FieldString("handler", handler.Name()),
		FieldString("trigger", trigger),
		FieldString("guild", msg.GuildID),
		FieldString("channel", msg.ChannelID),
		FieldString("author", msg.Author.ID),
		FieldString("message", msg.ID),
	}

	if err != nil {
		// Log error
		bot.LogError("Handler error", logFields...)
	} else {
		// Log every trigger
		bot.LogInfo("Handler trigger", logFields...)
	}

	bot.EmitEvent(EventMessageReceived, handleInfo)
}
