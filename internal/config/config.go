package config

import (
	"errors"
	"os"
)

type Config struct {
	BotToken string
}

func Load() (*Config, error) {
	c := &Config{}

	var found bool
	if c.BotToken, found = os.LookupEnv("BOT_TOKEN"); !found {
		return nil, errors.New("BOT_TOKEN not specified")
	}

	return c, nil
}
