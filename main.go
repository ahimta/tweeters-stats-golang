package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Ahimta/tweeters-stats-golang/auth"
	"github.com/Ahimta/tweeters-stats-golang/config"
	"github.com/Ahimta/tweeters-stats-golang/handlers"
	"github.com/Ahimta/tweeters-stats-golang/middleware"
	"github.com/Ahimta/tweeters-stats-golang/services"
	"github.com/Ahimta/tweeters-stats-golang/usecases"
	newrelic "github.com/newrelic/go-agent"
)

func main() {
	c, err := config.New(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("CALLBACK_URL"),
		os.Getenv("PORT"),
		os.Getenv("HOMEPAGE"),
		os.Getenv("HOST"),
		os.Getenv("PROTOCOL"),
		os.Getenv("CORS_DOMAIN"),
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var app newrelic.Application
	if licenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY"); licenseKey != "" {
		config := newrelic.NewConfig("tweeters-stats-golang", licenseKey)
		app, err = newrelic.NewApplication(config)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
	}

	oauthClient, err := auth.NewOauth1Client(
		c.ConsumerKey,
		c.ConsumerSecret,
		c.CallbackURL,
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	tweetsService := services.NewTweetsService(oauthClient)
	mux := http.NewServeMux()

	route(mux, app, "/", handlers.Homepage())
	route(mux, app, "/login/twitter", handlers.Login(usecases.Login, oauthClient))
	route(
		mux,
		app,
		"/oauth/twitter/callback",
		handlers.OauthTwitter(
			usecases.Oauth1Callback,
			c,
			oauthClient,
		),
	)
	route(mux, app, "/logout", handlers.Logout())
	route(
		mux,
		app,
		"/tweeters-stats",
		handlers.TweetersStats(
			usecases.TweetersStats,
			tweetsService,
		),
	)

	fmt.Printf("Server running on %s://%s\n", c.Protocol, c.Host)
	http.ListenAndServe(
		fmt.Sprintf(":%s", c.Port),
		middleware.Apply(mux, os.Stdout, c),
	)
}

func route(
	mux *http.ServeMux,
	app newrelic.Application,
	path string,
	handler http.HandlerFunc,
) {
	if app == nil {
		mux.HandleFunc(path, handler)
		return
	}

	mux.HandleFunc(newrelic.WrapHandleFunc(app, path, handler))
}
