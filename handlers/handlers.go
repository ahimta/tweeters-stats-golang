package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ahimta/tweeters-stats-golang/auth"
	"github.com/Ahimta/tweeters-stats-golang/config"
	"github.com/Ahimta/tweeters-stats-golang/entities"
	"github.com/Ahimta/tweeters-stats-golang/services"
	"github.com/Ahimta/tweeters-stats-golang/usecases"
)

// LoginHandlerFactory blablabla
func LoginHandlerFactory(c *config.Config, oauthClient auth.Oauth1Client) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		oauthLoginResult, err := usecases.Login(oauthClient)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			w.Write([]byte("500 - oops"))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "oauthRequestSecret",
			Value: oauthLoginResult.RequestSecret,
			Path:  "/",
		})

		http.Redirect(w, r, oauthLoginResult.AuthorizationURL.String(), http.StatusFound)
	}
}

// OauthTwitterHandlerFactory blablabla
func OauthTwitterHandlerFactory(c *config.Config, oauthClient auth.Oauth1Client) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		requestSecret := getCookieValue(r, "oauthRequestSecret")
		handleOauthResult, err := usecases.HandleOauth1Callback(oauthClient, requestSecret, r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			w.Write([]byte("500 - oops"))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "accessToken",
			Value: handleOauthResult.AccessToken,
			Path:  "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:  "accessSecret",
			Value: handleOauthResult.AccessSecret,
			Path:  "/",
		})

		http.Redirect(w, r, "/tweeters-stats", http.StatusFound)
	}
}

type getTweetersStatsResponse struct {
	Data []*entities.TweeterStats `json:"data"`
}

// GetTweetersStatsHandlerFactory blablabla
func GetTweetersStatsHandlerFactory(c *config.Config, oauthClient auth.Oauth1Client) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := getCookieValue(r, "accessToken")
		accessSecret := getCookieValue(r, "accessSecret")
		stats, err := usecases.GetTweetersStats(services.NewTweetsService(oauthClient), accessToken, accessSecret)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			w.Write([]byte("500 - oops"))
			return
		}

		json.NewEncoder(w).Encode(&getTweetersStatsResponse{stats})
	}
}

func getCookieValue(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)

	if err != nil {
		return ""
	}

	return cookie.Value
}
