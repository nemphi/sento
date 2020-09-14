package discord

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nemphi/sento"

	"github.com/bwmarrin/discordgo"
)

// Bot is a sento-powered bot application
type Bot struct {
	sento.Bot

	Sess     *discordgo.Session
	handlers map[string]Handler
	cfg      *Config
}

// New returns a new sento-powered discord bot
func New(options ...Option) (bot *Bot, err error) {
	bot = &Bot{
		Bot: sento.New(),
	}
	for _, op := range options {
		err = op(bot)
		if err != nil {
			break
		}
	}
	return
}

// SetConfig of a bot
func (bot *Bot) SetConfig(cfg *Config) {
	bot.cfg = cfg
}

// Start an instance of the bot
func (bot *Bot) Start() (err error) {
	bot.LogInfo("Creating session")
	bot.Sess, err = discordgo.New("Bot " + bot.cfg.Token)
	if err != nil {
		// TODO: Maybe modify error message
		// Could not connect to host/discord
		bot.LogInfo("Error creating session")
		return err
	}

	// Add handlers
	for _, handler := range bot.handlers {
		err = handler.Start(bot)
		if err != nil {
			// TODO: Maybe modify error message
			// Error while starting handler
			return err
		}
	}
	bot.Sess.AddHandler(bot.handleCreateMessage)

	bot.LogInfo("Opening the connection")
	err = bot.Sess.Open()
	if err != nil {
		bot.LogInfo("Error opening the connection")
		return
	}

	bot.LogInfo("Listening . . .")

	return
}

// Stop the bot
func (bot *Bot) Stop() (err error) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	for _, handler := range bot.handlers {
		err = handler.Stop(bot)
		if err != nil {
			// TODO: Maybe modify error message
			// Error while stoping handler
			return err
		}
	}
	err = bot.Sess.Close()
	return
}

func (bot *Bot) handleCreateMessage(sess *discordgo.Session, msg *discordgo.MessageCreate) {
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
	handler.Handle(bot, HandleInfo{
		AuthorID:    msg.Author.ID,
		ChannelID:   msg.ChannelID,
		Trigger:     trigger,
		Timestamp:   time.Now(),
		FullMessage: msg.Message,
	})
}
