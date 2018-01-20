package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ahimta/tweeters-stats-golang/auth"
	"github.com/Ahimta/tweeters-stats-golang/entities"
	"github.com/Ahimta/tweeters-stats-golang/services"
	"github.com/Ahimta/tweeters-stats-golang/usecases"
)

type loginUsecaseFunc func(client auth.Oauth1Client) (*usecases.OauthLoginResult, error)
type handleOauth1CallbackUsecaseFunc func(oauthClient auth.Oauth1Client, requestSecret string, r *http.Request) (*usecases.HandleOauth1CallbackResult, error)
type getTweetersStatsUsecaseFunc func(tweetsService services.TweetsService, accessToken, accessSecret string) ([]*entities.TweeterStats, error)

// LoginHandlerFactory blablabla
func LoginHandlerFactory(loginUsecase loginUsecaseFunc, oauthClient auth.Oauth1Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		oauthLoginResult, err := loginUsecase(oauthClient)

		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/", http.StatusFound)
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

// LogoutHandlerFactory blablabla
func LogoutHandlerFactory() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
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

// OauthTwitterHandlerFactory blablabla
func OauthTwitterHandlerFactory(handleOauth1CallbackUsecase handleOauth1CallbackUsecaseFunc, oauthClient auth.Oauth1Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		requestSecret := getCookieValue(r, "oauthRequestSecret")
		handleOauthResult, err := handleOauth1CallbackUsecase(oauthClient, requestSecret, r)

		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/", http.StatusFound)
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

		http.Redirect(w, r, "http://127.0.0.1:8000/Main.elm", http.StatusFound)
	}
}

// GetTweetersStatsResponse blablabla
type GetTweetersStatsResponse struct {
	Data []*entities.TweeterStats `json:"data"`
}

// GetTweetersStatsHandlerFactory blablabla
func GetTweetersStatsHandlerFactory(getTweetersStatsUsecase getTweetersStatsUsecaseFunc, tweetsService services.TweetsService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := getCookieValue(r, "accessToken")
		accessSecret := getCookieValue(r, "accessSecret")
		stats, err := getTweetersStatsUsecase(tweetsService, accessToken, accessSecret)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(&GetTweetersStatsResponse{stats})
	}
}

func getCookieValue(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)

	if err != nil {
		return ""
	}

	return cookie.Value
}
