package auth

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/dghubble/oauth1"
)

// Oauth1Client blablabla
type Oauth1Client interface {
	AccessToken(requestToken, requestSecret, verifier string) (accessToken, accessSecret string, err error)
	AuthorizationURL(requestToken string) (*url.URL, error)
	HTTPClient(accessToken, accessSecret string) (*http.Client, error)
	RequestToken() (requestToken, requestSecret string, err error)

	ParseAuthorizationCallback(r *http.Request) (requestToken, verifier string, err error)
}

// oauth1Client blabla
type oauth1Client struct {
	accessTokenImpl                func(requestToken, requestSecret, verifier string) (accessToken, accessSecret string, err error)
	authorizationURLImpl           func(requestToken string) (*url.URL, error)
	clientImpl                     func(ctx context.Context, t *oauth1.Token) *http.Client
	newTokenImpl                   func(token, tokenSecret string) *oauth1.Token
	parseAuthorizationCallbackImpl func(r *http.Request) (requestToken, verifier string, err error)
	requestTokenImpl               func() (requestToken, requestSecret string, err error)
}

// NewOauth1Client blabla
func NewOauth1Client(consumerKey, consumerSecret, callbackURL string) (Oauth1Client, error) {
	if consumerKey == "" || consumerSecret == "" || callbackURL == "" {
		return nil, errors.New("auth: consumerKey, consumerSecret, or callbackURL is missing -_-")
	}

	config := &oauth1.Config{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		CallbackURL:    callbackURL,
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: "https://api.twitter.com/oauth/request_token",
			AuthorizeURL:    "https://api.twitter.com/oauth/authorize",
			AccessTokenURL:  "https://api.twitter.com/oauth/access_token",
		},
	}

	return &oauth1Client{
		accessTokenImpl:                config.AccessToken,
		authorizationURLImpl:           config.AuthorizationURL,
		clientImpl:                     config.Client,
		newTokenImpl:                   oauth1.NewToken,
		parseAuthorizationCallbackImpl: oauth1.ParseAuthorizationCallback,
		requestTokenImpl:               config.RequestToken,
	}, nil
}

// AccessToken blabla
func (client *oauth1Client) AccessToken(requestToken, requestSecret, verifier string) (accessToken, accessSecret string, err error) {
	if requestToken == "" || requestSecret == "" || verifier == "" {
		return "", "", errors.New("auth: missing requestToken, requestSecret, or verifier -_-")
	}

	return client.accessTokenImpl(requestToken, requestSecret, verifier)
}

// AuthorizationURL blabla
func (client *oauth1Client) AuthorizationURL(requestToken string) (*url.URL, error) {
	if requestToken == "" {
		return nil, errors.New("auth: missing requestToken")
	}

	return client.authorizationURLImpl(requestToken)
}

// HTTPClient blabla
func (client *oauth1Client) HTTPClient(accessToken, accessSecret string) (*http.Client, error) {
	if accessToken == "" || accessSecret == "" {
		return nil, errors.New("auth: missing accessToken or accessSecret")
	}

	token := client.newTokenImpl(accessToken, accessSecret)
	return client.clientImpl(oauth1.NoContext, token), nil
}

// RequestToken blabla
func (client *oauth1Client) RequestToken() (requestToken, requestSecret string, err error) {
	return client.requestTokenImpl()
}

// ParseAuthorizationCallback blabla
func (client *oauth1Client) ParseAuthorizationCallback(r *http.Request) (requestToken, verifier string, err error) {
	if r == nil {
		return "", "", errors.New("auth: missing request")
	}

	return client.parseAuthorizationCallbackImpl(r)
}
