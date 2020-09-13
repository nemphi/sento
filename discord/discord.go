package discord

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Bot is a sento-powered bot application
type Bot struct {
	Sess     *discordgo.Session
	handlers map[string]Handler
	cfg      *Config
}

// New returns a new sento-powered discord bot
func New(options ...Option) (bot *Bot, err error) {
	bot = &Bot{}
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
	bot.Sess, err = discordgo.New("Bot " + bot.cfg.Token)
	if err != nil {
		// TODO: Maybe modify error message
		// Could not connect to host/discord
		return err
	}

	// Add handlers
	bot.Sess.AddHandler(bot.handleCreateMessage)

	return
}

// Stop the bot
func (bot *Bot) Stop() (err error) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = bot.Sess.Close()
	return
}

func (bot *Bot) handleCreateMessage(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	// TODO: Fetch prefix from database

	// TODO: prefix check

	handler, exists := bot.handlers[msg.Content]
	if !exists {
		// TODO: handle case
	}

	// Handle message
	// TODO: Make async
	handler.Handle(bot, HandleInfo{
		AuthorID:    msg.Author.ID,
		ChannelID:   msg.ChannelID,
		Trigger:     msg.Content, // TODO: use trigger
		Timestamp:   time.Now(),
		FullMessage: msg.Message,
	})
}
