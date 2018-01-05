package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Ahimta/tweeters-stats-golang/auth"
	"github.com/Ahimta/tweeters-stats-golang/config"
	"github.com/Ahimta/tweeters-stats-golang/handlers"
	"github.com/Ahimta/tweeters-stats-golang/usecases"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func main() {
	c, err := config.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"), os.Getenv("CALLBACK_URL"), os.Getenv("PORT"))

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	oauthClient, err := auth.NewOauth1Client(c.ConsumerKey, c.ConsumerSecret, c.CallbackURL)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)
	mux.HandleFunc("/login/twitter", handlers.LoginHandlerFactory(usecases.Login, oauthClient))
	mux.HandleFunc("/oauth/twitter/callback", handlers.OauthTwitterHandlerFactory(usecases.HandleOauth1Callback, oauthClient))
	mux.HandleFunc("/tweeters-stats", handlers.GetTweetersStatsHandlerFactory(usecases.GetTweetersStats, oauthClient))

	fmt.Printf("Server running on http://localhost:%s\n", c.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", c.Port), mux)
}
