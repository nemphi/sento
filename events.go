package sento

// EventListener interface used for defining listeners
type EventListener interface {
	Type() EventType
	Handle(bot *Bot, data interface{})
}

// EventType indicates the supported event types
type EventType int

const (
	// EventConnected emitted when the discord session opens
	EventConnected EventType = iota
	// EventMessageSent emitted when the bot sends a message
	EventMessageSent
	// EventMessageReceived emitted when the bot processes a message
	EventMessageReceived
	// EventDisconnected emitted when the discord session closes
	EventDisconnected
)

// EmitEvent and broadcast it to all listeners
func (bot *Bot) EmitEvent(eventType EventType, data interface{}) {
	listeners, notEmpty := bot.listeners[eventType]
	if notEmpty {
		for _, listener := range listeners {
			go listener.Handle(bot, data)
		}
	}
}
