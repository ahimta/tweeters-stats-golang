package config

import (
	"os"
)

// Config blablabla
type Config struct {
	CallbackURL    string
	ConsumerKey    string
	ConsumerSecret string
	Port           string
}

var config *Config

// NewConfig blablabla
func NewConfig(callbackURL, consumerKey, consumerSecret, port string) *Config {
	return &Config{callbackURL, consumerKey, consumerSecret, port}
}

// GetConfig blablabla
func GetConfig() *Config {
	if config != nil {
		return config
	}

	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	callbackURL := os.Getenv("CALLBACK_URL")
	port := os.Getenv("PORT")

	if consumerKey == "" || consumerSecret == "" || callbackURL == "" {
		panic("Missing CONSUMER_KEY or CONSUMER_SECRET or CALLBACK_URL environment variables -_-!")
	}

	if port == "" {
		port = "8080"
	}

	config = NewConfig(callbackURL, consumerKey, consumerSecret, port)
	return config
}
