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
}

var config *Config

// NewConfig blablabla
func NewConfig(consumerKey, consumerSecret, callbackURL, port string) (
	*Config, error,
) {

	if consumerKey == "" ||
		consumerSecret == "" ||
		callbackURL == "" ||
		port == "" {
		return nil, errors.New("config: a required parameter is missing -_-")
	}

	return &Config{consumerKey, consumerSecret, callbackURL, port}, nil
}
