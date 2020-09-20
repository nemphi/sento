package sento

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Option for a sent-powered discord bot
type Option func(*Bot) error

// UseConfigFile makes a robot use the indicated
// config file. Note that the config MUST be in
// TOML format.
func UseConfigFile(path string) Option {
	return func(bot *Bot) error {
		cfg := &Config{}
		file, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		if _, err = toml.Decode(string(file), cfg); err != nil {
			return err
		}
		bot.cfg = cfg
		return nil
	}
}

// UseConfig sets the config for a bot
func UseConfig(cfg *Config) Option {
	return func(bot *Bot) error {
		bot.cfg = cfg
		return nil
	}
}

// UseListeners sets the listeners for a bot
func UseListeners(listeners ...EventListener) Option {
	return func(bot *Bot) error {
		if bot.listeners == nil {
			bot.listeners = make(map[EventType][]EventListener)
		}
		for _, listener := range listeners {
			bot.listeners[listener.Type()] = append(bot.listeners[listener.Type()], listener)
		}
		return nil
	}
}

// UseHandlers sets the handlers for a bot
func UseHandlers(handlers ...Handler) Option {
	return func(bot *Bot) error {
		if bot.handlers == nil {
			bot.handlers = make(map[string]Handler)
		}
		for _, handler := range handlers {
			for _, trigger := range handler.Triggers() {
				bot.handlers[trigger] = handler
			}
		}
		return nil
	}
}
