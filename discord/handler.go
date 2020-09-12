package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// Handler listens for an specific command and
// contain all the logic necesary for it
type Handler interface {
	// Name of the handler
	Name() string
	// Triggers for a given handler
	Triggers() string
	// Handle the trigger instance
	Handle(bot *Bot, info HandleInfo) error

	// Start runs when the bot connection has been made
	// and is adding all handlers
	Start()
	// Stop runs when the bot is being shut down
	Stop()
}

// HandleInfo about a single trigger instance
type HandleInfo struct {
	Trigger     string
	ChannelID   string
	AuthorID    string
	Message     string
	FullMessage *discordgo.Message
	Timestamp   time.Time
}

// HandleError ocurred while dealing with a trigger
type HandleError struct {
	HandlerName   string
	MessageID     string
	OriginalError error
}

func (he HandleError) Error() string {
	return "HandleError(" + he.HandlerName + ") " + he.OriginalError.Error()
}

// --------------- Just an example implementation -------------

// WILL be moved out of here, meanwhile it stays here

type pingPong struct {
}

func (p pingPong) Start() {}
func (p pingPong) Stop()  {}

func (p pingPong) Name() string {
	return "PingPong"
}

func (p pingPong) Triggers() []string {
	return []string{
		"*",
	}
}

func (p pingPong) Handle(bot *Bot, info HandleInfo) (err error) {
	if info.Message == "ping" {
		_, err = bot.Sess.ChannelMessageSend(info.ChannelID, "pong!")
	} else if info.Message == "pong" {
		_, err = bot.Sess.ChannelMessageSend(info.ChannelID, "ping!")
	}

	if err != nil {
		err = HandleError{
			HandlerName:   p.Name(),
			MessageID:     info.FullMessage.ID,
			OriginalError: err,
		}
	}
	return
}

// ------------ End of example ----------
