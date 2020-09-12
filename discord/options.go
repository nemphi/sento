package discord

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
		bot.SetConfig(cfg)
		return nil
	}
}
