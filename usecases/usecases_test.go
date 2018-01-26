package usecases

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/Ahimta/tweeters-stats-golang/entities"
)

type oauthClient struct {
	// AccessToken
	accessToken    string
	accessSecret   string
	accessTokenErr error

	// AuthorizationURL
	url                   *url.URL
	authorizationURLError error

	// HTTPClient
	client *http.Client

	// RequestToken
	requestToken      string
	requestSecret     string
	requestTokenError error

	// ParseAuthorizationCallback
	parseAuthorizationRequestToken  string
	verifier                        string
	parseAuthorizationCallbackError error
}

type tweetsService struct {
	tweeters []*entities.Tweeter
	err      error
}

func (service *tweetsService) Tweeters(
	accessToken,
	accessSecret string,
) (
	[]*entities.Tweeter, error,
) {

	return service.tweeters, service.err
}

func (client *oauthClient) AccessToken(
	requestToken,
	requestSecret,
	verifier string) (
	accessToken, accessSecret string, err error,
) {

	return client.accessToken, client.accessSecret, client.accessTokenErr
}

func (client *oauthClient) AuthorizationURL(requestToken string) (
	*url.URL, error,
) {

	return client.url, client.authorizationURLError
}

func (client *oauthClient) HTTPClient(accessToken, accessSecret string) (
	*http.Client, error,
) {

	return client.client, nil
}

func (client *oauthClient) RequestToken() (
	requestToken, requestSecret string, err error,
) {

	return client.requestToken, client.requestSecret, client.requestTokenError
}

func (client *oauthClient) ParseAuthorizationCallback(r *http.Request) (
	requestToken, verifier string, err error,
) {

	return client.parseAuthorizationRequestToken,
		client.verifier,
		client.parseAuthorizationCallbackError
}

func TestTweetersStats(t *testing.T) {
	//
	stats, err := TweetersStats(
		&tweetsService{
			tweeters: []*entities.Tweeter{
				&entities.Tweeter{
					FullName: "John Smith1",
					Username: "jsmith1",
				},
				&entities.Tweeter{
					FullName: "John Smith0",
					Username: "jsmith0",
				},
				&entities.Tweeter{
					FullName: "John Smith2",
					Username: "jsmith2",
				},
				&entities.Tweeter{
					FullName: "John Smith0",
					Username: "jsmith0",
				},
			},
			err: nil,
		},
		"blablabla",
		"blablabla",
	)

	if err != nil {
		t.Errorf(err.Error())
	}

	if stats[0].FullName != "John Smith0" ||
		stats[0].Username != "jsmith0" ||
		stats[0].TweetsCount != 2 {

		t.Errorf("Should return stats related to tweeters from TweetService")
	}

	if len(stats) != 3 {
		t.Errorf("Whaaaat!")
	}

	//
	stats, err = TweetersStats(
		&tweetsService{
			tweeters: nil,
			err:      nil,
		},
		"blablabla",
		"",
	)

	if err == nil {
		t.Errorf("Should return an error when a parameter is missing")
	}

	if len(stats) != 0 {
		t.Errorf("Whaaaat!")
	}

	//
	stats, err = TweetersStats(
		&tweetsService{
			tweeters: nil,
			err:      nil,
		},
		"blablabla",
		"",
	)

	if err == nil {
		t.Errorf("Should return an error when a parameter is missing")
	}

	if len(stats) != 0 {
		t.Errorf("Whaaaat!")
	}

	//
	stats, err = TweetersStats(
		&tweetsService{
			tweeters: nil,
			err:      errors.New("blablabla"),
		},
		"blablabla",
		"blabla",
	)

	if !(err != nil) {
		t.Errorf("Should return an error when TweetService returns an error")
	}

	if len(stats) != 0 {
		t.Errorf("Whaaaat!")
	}
}

func TestHandleOauth1Callback(t *testing.T) {
	//
	result, err := Oauth1Callback(
		&oauthClient{
			parseAuthorizationRequestToken:  "requestToken",
			verifier:                        "verifier",
			parseAuthorizationCallbackError: nil,

			accessToken:    "accessToken",
			accessSecret:   "accessSecret",
			accessTokenErr: nil,
		},
		"blablabla",
		&http.Request{},
	)

	if err != nil {
		t.Errorf(err.Error())
	}

	if result.AccessToken != "accessToken" ||
		result.AccessSecret != "accessSecret" {

		t.Errorf("Should match Oauth1Client.AccessToken return value")
	}

	//
	result, err = Oauth1Callback(&oauthClient{}, "blablabla", nil)

	if err == nil || result != nil {
		t.Errorf("Whaaat!")
	}

	//
	result, err = Oauth1Callback(
		&oauthClient{
			parseAuthorizationCallbackError: errors.New("blablabla"),
		},
		"blablabla",
		&http.Request{},
	)

	if err == nil || result != nil {
		t.Errorf("Whaaat!")
	}

	//
	result, err = Oauth1Callback(
		&oauthClient{
			accessTokenErr: errors.New("blablabla"),
		},
		"blablabla",
		&http.Request{},
	)

	if err == nil || result != nil {
		t.Errorf("Whaaat!")
	}
}

func TestLogin(t *testing.T) {
	//
	result, err := Login(&oauthClient{
		requestToken:  "requestToken",
		requestSecret: "requestSecret",
		url:           &url.URL{Path: "blablabla"},
	})

	if err != nil ||
		result.RequestSecret != "requestSecret" ||
		result.AuthorizationURL.Path != "blablabla" {

		t.Errorf(
			"Should return results from client.{RequestToken(),AuthorizationURL()}",
		)
	}

	//
	result, err = Login(&oauthClient{requestTokenError: errors.New("blablabla")})

	if err == nil || result != nil {
		t.Errorf("Whaaat!")
	}
}
