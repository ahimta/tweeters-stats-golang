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
)

func main() {
	c, err := config.New(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("CALLBACK_URL"),
		os.Getenv("PORT"),
		os.Getenv("HOMEPAGE"),
		os.Getenv("CORS_DOMAIN"),
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
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

	mux.HandleFunc("/", handlers.Homepage())
	mux.HandleFunc("/login/twitter", handlers.Login(usecases.Login, oauthClient))
	mux.HandleFunc(
		"/oauth/twitter/callback",
		handlers.OauthTwitter(
			usecases.Oauth1Callback,
			c,
			oauthClient,
		),
	)
	mux.HandleFunc("/logout", handlers.Logout())
	mux.HandleFunc(
		"/tweeters-stats",
		handlers.TweetersStats(
			usecases.TweetersStats,
			tweetsService,
		),
	)

	fmt.Printf("Server running on http://localhost:%s\n", c.Port)
	http.ListenAndServe(
		fmt.Sprintf(":%s", c.Port),
		middleware.Apply(mux, os.Stdout, c),
	)
}
