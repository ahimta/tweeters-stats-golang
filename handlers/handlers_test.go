package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/Ahimta/tweeters-stats-golang/auth"
	"github.com/Ahimta/tweeters-stats-golang/config"
	"github.com/Ahimta/tweeters-stats-golang/entities"
	"github.com/Ahimta/tweeters-stats-golang/services"
	"github.com/Ahimta/tweeters-stats-golang/usecases"
)

func TestLogin(t *testing.T) {
	t.Run(
		"should use underlying implementation and redirect to correct URL",
		func(t *testing.T) {

			req, err := http.NewRequest("GET", "/login/twitter", nil)

			if err != nil {
				t.Fatal(err)
			}

			authorizationURL := &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "authorizationUrl",
			}

			oauthClient, err := auth.NewOauth1Client(
				"consumerKey",
				"consumerSecret",
				"callbackURL",
			)

			if err != nil {
				t.Fatal(err)
			}

			usecase := func(client auth.Oauth1Client) (*usecases.LoginResult, error) {
				if client != oauthClient {
					t.Errorf("oauthClient not passed to login usecase -_-")
				}

				return &usecases.LoginResult{
					AuthorizationURL: authorizationURL,
					RequestSecret:    "requestSecret",
				}, nil
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Login(usecase, oauthClient))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusFound {
				t.Errorf("Expected 302 HTTP status code")
			}

			setCookie := rr.Header().Get("Set-Cookie")
			if setCookie != "oauthRequestSecret=requestSecret; Path=/" {
				t.Errorf("Incorrect Set-Cookie value: %v", setCookie)
			}

			location := rr.Header().Get("Location")
			if location != authorizationURL.String() {
				t.Errorf(
					"Incorrect Location value: %v != %v",
					location,
					authorizationURL.String(),
				)
			}
		})

	t.Run(
		"should use return an error when the underlying implementation does",
		func(t *testing.T) {

			req, err := http.NewRequest("GET", "/login/twitter", nil)

			if err != nil {
				t.Fatal(err)
			}

			if err != nil {
				t.Fatal(err)
			}

			usecase := func(client auth.Oauth1Client) (*usecases.LoginResult, error) {
				return nil, errors.New("whaaat -_-")
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Login(usecase, nil))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusFound {
				t.Errorf("Expected 302 HTTP status code")
			}

			if setCookie := rr.Header().Get("Set-Cookie"); setCookie != "" {
				t.Errorf("Incorrect Set-Cookie value: %v", setCookie)
			}

			if location := rr.Header().Get("Location"); location != "/" {
				t.Errorf("Incorrect Location value: %v != %v", location, "/")
			}
		})
}

func TestOauthTwitter(t *testing.T) {
	c, err := config.New(
		"consumerKey",
		"consumerSecret",
		"callbackURL",
		"80",
		"/",
		"localhost",
		"http",
		"",
	)

	if err != nil {
		t.Errorf(err.Error())
	}

	t.Run(
		"should use underlying implementation and redirect to correct URL",
		func(t *testing.T) {

			req, err := http.NewRequest("GET", "/oauth/twitter/callback", nil)
			req.AddCookie(&http.Cookie{
				Name:  "oauthRequestSecret",
				Value: "requestSecret",
			})

			if err != nil {
				t.Fatal(err)
			}

			result := &usecases.Oauth1CallbackResult{
				AccessToken:  "accessToken",
				AccessSecret: "accessSecret",
			}

			oauthClient, err := auth.NewOauth1Client(
				"consumerKey",
				"consumerSecret",
				"callbackURL",
			)

			if err != nil {
				t.Fatal(err)
			}

			usecase := func(
				client auth.Oauth1Client,
				requestSecret string,
				r *http.Request) (
				*usecases.Oauth1CallbackResult, error,
			) {

				if client != oauthClient ||
					requestSecret != "requestSecret" ||
					r == nil {

					t.Errorf("parameters not passed to oauth usecase correctly -_-")
				}

				return result, nil
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(OauthTwitter(usecase, c, oauthClient))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusFound {
				t.Errorf("Expected 302 HTTP status code")
			}

			setCookie := rr.Header().Get("Set-Cookie")
			if setCookie != "accessToken=accessToken; Path=/" {
				t.Errorf("Incorrect Set-Cookie value: %v", setCookie)
			}

			location := rr.Header().Get("Location")
			if location != c.Homepage {
				t.Errorf(
					"Incorrect Location value: %v != %v",
					location,
					c.Homepage,
				)
			}
		})

	t.Run(
		"should use underlying implementation and redirect to correct URL",
		func(t *testing.T) {

			req, err := http.NewRequest("GET", "/oauth/twitter/callback", nil)

			if err != nil {
				t.Fatal(err)
			}

			usecase := func(
				client auth.Oauth1Client,
				requestSecret string,
				r *http.Request) (
				*usecases.Oauth1CallbackResult, error,
			) {

				return nil, errors.New("whaaat -_-")
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(OauthTwitter(usecase, c, nil))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusFound {
				t.Errorf("Expected 302 HTTP status code")
			}

			if setCookie := rr.Header().Get("Set-Cookie"); setCookie != "" {
				t.Errorf("Incorrect Set-Cookie value: %v", setCookie)
			}

			if location := rr.Header().Get("Location"); location != c.Homepage {
				t.Errorf("Incorrect Location value: %v != %v", location, c.Homepage)
			}
		})
}

func TestTweetersStats(t *testing.T) {
	t.Run(
		"should use underlying implementation and redirect to correct URL",
		func(t *testing.T) {

			req, err := http.NewRequest("GET", "/tweeters-stats", nil)
			req.AddCookie(&http.Cookie{
				Name:  "accessToken",
				Value: "accessToken",
			})
			req.AddCookie(&http.Cookie{
				Name:  "accessSecret",
				Value: "accessSecret",
			})

			if err != nil {
				t.Fatal(err)
			}

			result := []*entities.TweeterStats{
				&entities.TweeterStats{
					FullName:    "John Smith",
					Username:    "jsmith",
					TweetsCount: 3,
				},
			}

			oauthClient, err := auth.NewOauth1Client(
				"consumerKey",
				"consumerSecret",
				"callbackURL",
			)

			if err != nil {
				t.Fatal(err)
			}

			tweetsService := services.NewTweetsService(oauthClient)

			usecase := func(
				service services.TweetsService, accessToken,
				accessSecret string,
			) (
				[]*entities.TweeterStats, error,
			) {

				if service != tweetsService ||
					accessToken != "accessToken" ||
					accessSecret != "accessSecret" {

					t.Errorf(
						"parameters not passed to tweeters-stats usecase correctly -_-",
					)
				}

				return result, nil
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(TweetersStats(usecase, tweetsService))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Expected 200 HTTP status code")
			}

			if setCookie := rr.Header().Get("Set-Cookie"); setCookie != "" {
				t.Errorf("Incorrect Set-Cookie value: %v", setCookie)
			}

			var responseBody TweetersStatsResponse
			json.NewDecoder(rr.Body).Decode(&responseBody)

			if !reflect.DeepEqual(responseBody.Data, result) {
				t.Errorf(
					"Incorrect response body: %v, expected: %v",
					responseBody.Data,
					result,
				)
			}
		})

	t.Run("should handle the usecase error with a 401 code", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/tweeters-stats", nil)

		if err != nil {
			t.Fatal(err)
		}

		oauthClient, err := auth.NewOauth1Client(
			"consumerKey",
			"consumerSecret",
			"callbackURL",
		)

		if err != nil {
			t.Fatal(err)
		}

		tweetsService := services.NewTweetsService(oauthClient)

		usecase := func(
			service services.TweetsService, accessToken,
			accessSecret string,
		) (
			[]*entities.TweeterStats, error,
		) {

			return nil, errors.New("whaaat -_-")
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(TweetersStats(usecase, tweetsService))
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("Expected 401 HTTP status code")
		}

		if setCookie := rr.Header().Get("Set-Cookie"); setCookie != "" {
			t.Errorf("Incorrect Set-Cookie value: %v", setCookie)
		}
	})
}
