package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Ahimta/tweeters-stats-golang/auth"
	"github.com/Ahimta/tweeters-stats-golang/config"
	"github.com/Ahimta/tweeters-stats-golang/entities"
	"github.com/Ahimta/tweeters-stats-golang/services"
	"github.com/Ahimta/tweeters-stats-golang/usecases"
)

type loginUsecaseFunc func(client auth.Oauth1Client) (
	*usecases.LoginResult, error,
)

type oauth1CallbackUsecaseFunc func(
	oauthClient auth.Oauth1Client,
	requestSecret string,
	r *http.Request) (
	*usecases.Oauth1CallbackResult, error,
)

type tweetersStatsUsecaseFunc func(
	tweetsService services.TweetsService, accessToken,
	accessSecret string,
) (
	[]*entities.TweeterStats, error,
)

// Homepage blablabla
func Homepage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("index.html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(data)
	}
}

// Login blablabla
func Login(
	usecase loginUsecaseFunc,
	client auth.Oauth1Client,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		result, err := usecase(client)

		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "oauthRequestSecret",
			Value: result.RequestSecret,
			Path:  "/",
		})

		http.Redirect(w, r, result.AuthorizationURL.String(), http.StatusFound)
	}
}

// Logout blablabla
func Logout() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "oauthRequestSecret",
			Value: "",
			Path:  "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "accessToken",
			Value: "",
			Path:  "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "accessSecret",
			Value: "",
			Path:  "/",
		})

		w.WriteHeader(http.StatusNoContent)
	}
}

// OauthTwitter blablabla
func OauthTwitter(
	usecase oauth1CallbackUsecaseFunc,
	c *config.Config,
	client auth.Oauth1Client,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		requestSecret := cookieValue(r, "oauthRequestSecret")
		result, err := usecase(client, requestSecret, r)

		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, c.Homepage, http.StatusFound)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "accessToken",
			Value: result.AccessToken,
			Path:  "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:  "accessSecret",
			Value: result.AccessSecret,
			Path:  "/",
		})

		http.Redirect(w, r, c.Homepage, http.StatusFound)
	}
}

// TweetersStatsResponse blablabla
type TweetersStatsResponse struct {
	Data []*entities.TweeterStats `json:"data"`
}

// TweetersStats blablabla
func TweetersStats(
	usecase tweetersStatsUsecaseFunc,
	service services.TweetsService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		accessToken := cookieValue(r, "accessToken")
		accessSecret := cookieValue(r, "accessSecret")
		stats, err := usecase(service, accessToken, accessSecret)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(&TweetersStatsResponse{stats})
	}
}

func cookieValue(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)

	if err != nil {
		return ""
	}

	return cookie.Value
}
