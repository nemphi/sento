package sento

import (
	"io/ioutil"
	"sync"

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
			bot.listeners = make(map[EventType]chan EventData)
		}
		for _, listener := range listeners {
			listenerChan, exists := bot.listeners[listener.Type()]
			if exists {
				bot.listeners[listener.Type()] = merge(listenerChan, listener.Chan())
			} else {
				listenerChan := listener.Chan()
				bot.listeners[listener.Type()] = listenerChan
			}
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

func merge(cs ...chan EventData) chan EventData {
	out := make(chan EventData)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan EventData) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
