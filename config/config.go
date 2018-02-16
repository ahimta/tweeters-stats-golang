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

	Host     string
	Protocol string

	CorsDomain string
}

// New blablabla
func New(
	consumerKey,
	consumerSecret,
	callbackURL,
	port,
	homepage,
	host,
	protocol,
	corsDomain string) (
	*Config, error,
) {

	if consumerKey == "" ||
		consumerSecret == "" ||
		callbackURL == "" ||
		port == "" ||
		homepage == "" ||
		host == "" ||
		protocol == "" {
		return nil, errors.New("config: a required parameter is missing -_-")
	}

	return &Config{
		consumerKey,
		consumerSecret,
		callbackURL,
		port,
		homepage,

		host,
		protocol,

		corsDomain,
	}, nil
}
