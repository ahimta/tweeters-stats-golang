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

func (_tweetsService *tweetsService) FetchTweeters(accessToken, accessSecret string) ([]*entities.Tweeter, error) {
	return _tweetsService.tweeters, _tweetsService.err
}

func (_oauthClient *oauthClient) AccessToken(requestToken, requestSecret, verifier string) (accessToken, accessSecret string, err error) {
	return _oauthClient.accessToken, _oauthClient.accessSecret, _oauthClient.accessTokenErr
}

func (_oauthClient *oauthClient) AuthorizationURL(requestToken string) (*url.URL, error) {
	return _oauthClient.url, _oauthClient.authorizationURLError
}

func (_oauthClient *oauthClient) HTTPClient(accessToken, accessSecret string) (*http.Client, error) {
	return _oauthClient.client, nil
}

func (_oauthClient *oauthClient) RequestToken() (requestToken, requestSecret string, err error) {
	return _oauthClient.requestToken, _oauthClient.requestSecret, _oauthClient.requestTokenError
}

func (_oauthClient *oauthClient) ParseAuthorizationCallback(r *http.Request) (requestToken, verifier string, err error) {
	return _oauthClient.parseAuthorizationRequestToken, _oauthClient.verifier, _oauthClient.parseAuthorizationCallbackError
}

func TestGetTweetersStats(t *testing.T) {
	//
	stats, err := GetTweetersStats(&tweetsService{
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
	}, "blablabla", "blablabla")

	if err != nil {
		t.Errorf(err.Error())
	}

	if stats[0].FullName != "John Smith0" || stats[0].Username != "jsmith0" || stats[0].TweetsCount != 2 {
		t.Errorf("Should return stats related to tweeters from TweetService")
	}

	if len(stats) != 3 {
		t.Errorf("Whaaaat!")
	}

	//
	stats, err = GetTweetersStats(&tweetsService{
		tweeters: nil,
		err:      nil,
	}, "blablabla", "")

	if err == nil {
		t.Errorf("Should return an error when either accessToken or accessSecret is missing")
	}

	if len(stats) != 0 {
		t.Errorf("Whaaaat!")
	}

	//
	stats, err = GetTweetersStats(&tweetsService{
		tweeters: nil,
		err:      nil,
	}, "blablabla", "")

	if err == nil {
		t.Errorf("Should return an error when either accessToken or accessSecret is missing")
	}

	if len(stats) != 0 {
		t.Errorf("Whaaaat!")
	}

	//
	stats, err = GetTweetersStats(&tweetsService{
		tweeters: nil,
		err:      errors.New("blablabla"),
	}, "blablabla", "blabla")

	if !(err != nil) {
		t.Errorf("Should return an error when TweetService returns an error")
	}

	if len(stats) != 0 {
		t.Errorf("Whaaaat!")
	}
}

func TestHandleOauth1Callback(t *testing.T) {
	//
	result, err := HandleOauth1Callback(&oauthClient{
		parseAuthorizationRequestToken:  "requestToken",
		verifier:                        "verifier",
		parseAuthorizationCallbackError: nil,

		accessToken:    "accessToken",
		accessSecret:   "accessSecret",
		accessTokenErr: nil,
	}, "blablabla", &http.Request{})

	if err != nil {
		t.Errorf(err.Error())
	}

	if result.AccessToken != "accessToken" || result.AccessSecret != "accessSecret" {
		t.Errorf("Should match Oauth1Client.AccessToken return value")
	}

	//
	result, err = HandleOauth1Callback(&oauthClient{}, "blablabla", nil)

	if err == nil || result != nil {
		t.Errorf("Whaaat!")
	}

	//
	result, err = HandleOauth1Callback(&oauthClient{
		parseAuthorizationCallbackError: errors.New("blablabla"),
	}, "blablabla", &http.Request{})

	if err == nil || result != nil {
		t.Errorf("Whaaat!")
	}

	//
	result, err = HandleOauth1Callback(&oauthClient{
		accessTokenErr: errors.New("blablabla"),
	}, "blablabla", &http.Request{})

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

	if err != nil || result.RequestSecret != "requestSecret" || result.AuthorizationURL.Path != "blablabla" {
		t.Errorf("Should return results from client.{RequestToken(),AuthorizationURL()}")
	}

	//
	result, err = Login(&oauthClient{requestTokenError: errors.New("blablabla")})

	if err == nil || result != nil {
		t.Errorf("Whaaat!")
	}
}
