package sento

import (
	"log"
	"os"
)

// Bot is a sento-powered bot
type Bot struct {
	logger *log.Logger
}

// New sento-powered bot
func New() Bot {
	return Bot{
		logger: log.New(os.Stdout, "[SentoBot]", log.LstdFlags),
	}
}

// LogInfo logs the `msg` to the console
func (b Bot) LogInfo(msg string) {
	b.logger.Println(msg)
}
