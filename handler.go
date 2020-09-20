package sento

import (
	"github.com/bwmarrin/discordgo"
)

// Handler listens for an specific command and
// contain all the logic necesary for it
type Handler interface {
	// Name of the handler
	Name() string
	// Triggers for a given handler
	Triggers() []string
	// Handle the trigger instance
	Handle(bot *Bot, info HandleInfo) error

	// Start runs when the bot connection has been made
	// and is adding all handlers
	Start(bot *Bot) error
	// Stop runs when the bot is being shut down
	Stop(bot *Bot) error
}

// HandleInfo about a single trigger instance
type HandleInfo struct {
	Trigger        string
	GuildID        string
	ChannelID      string
	MessageID      string
	AuthorID       string
	MessageContent string
}

// Message returns the discord message that originated this trigger
func (hi HandleInfo) Message(bot *Bot) (*discordgo.Message, error) {
	return bot.Sess.ChannelMessage(hi.ChannelID, hi.MessageID)
}

// --------------- Just an example implementation -------------

// WILL be moved out of here, meanwhile it stays here

type pingPong struct {
}

func (p pingPong) Start(_ *Bot) (err error) { return }
func (p pingPong) Stop(_ *Bot) (err error)  { return }

func (p pingPong) Name() string {
	return "PingPong"
}

func (p pingPong) Triggers() []string {
	return []string{
		"ping",
		"pong",
	}
}

func (p pingPong) Handle(bot *Bot, info HandleInfo) (err error) {
	if info.Trigger == "ping" {
		err = bot.Send(info, "pong!")
	} else if info.Trigger == "pong" {
		err = bot.Send(info, "ping!")
	}

	return
}

// ------------ End of example ----------
