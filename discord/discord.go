package discord

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Bot is a sento-powered bot application
type Bot struct {
	sess     *discordgo.Session
	handlers []Handler
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
	bot.sess, err = discordgo.New("Bot " + "" /*TODO: Config bot token*/)
	if err != nil {
		// TODO: Maybe modify error message
		// Could not connect to host/discord
		return err
	}

	// Add handlers
	for _, handler := range bot.handlers {
		bot.sess.AddHandler(handler)
	}

	return
}

// Stop the bot
func (bot *Bot) Stop() (err error) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = bot.sess.Close()
	return
}
