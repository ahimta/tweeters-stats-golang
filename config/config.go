package config

import (
	"errors"
)

// Config blablabla
type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	CallbackURL    string
	Port           string
	Homepage       string
	CorsDomain     string
}

var config *Config

// NewConfig blablabla
func NewConfig(
	consumerKey,
	consumerSecret,
	callbackURL,
	port,
	homepage,
	corsDomain string) (
	*Config, error,
) {

	if consumerKey == "" ||
		consumerSecret == "" ||
		callbackURL == "" ||
		port == "" ||
		homepage == "" {
		return nil, errors.New("config: a required parameter is missing -_-")
	}

	return &Config{
		consumerKey,
		consumerSecret,
		callbackURL,
		port,
		homepage,
		corsDomain,
	}, nil
}
