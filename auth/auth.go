package auth

import (
	"net/http"
	"net/url"

	"github.com/Ahimta/tweeters-stats-golang/config"
	"github.com/dghubble/oauth1"
)

// Oauth1Client blablabla
type Oauth1Client interface {
	AccessToken(requestToken, requestSecret, verifier string) (accessToken, accessSecret string, err error)
	AuthorizationURL(requestToken string) (*url.URL, error)
	HTTPClient(accessToken, accessSecret string) *http.Client
	RequestToken() (requestToken, requestSecret string, err error)

	ParseAuthorizationCallback(r *http.Request) (requestToken, verifier string, err error)
}

// oauth1Client blabla
type oauth1Client struct {
	impl *oauth1.Config
}

// NewOauth1Client blabla
func NewOauth1Client(c *config.Config) Oauth1Client {
	config := &oauth1.Config{
		ConsumerKey:    c.ConsumerKey,
		ConsumerSecret: c.ConsumerSecret,
		CallbackURL:    c.CallbackURL,
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: "https://api.twitter.com/oauth/request_token",
			AuthorizeURL:    "https://api.twitter.com/oauth/authorize",
			AccessTokenURL:  "https://api.twitter.com/oauth/access_token",
		},
	}

	return &oauth1Client{config}
}

// AccessToken blabla
func (client *oauth1Client) AccessToken(requestToken, requestSecret, verifier string) (accessToken, accessSecret string, err error) {
	return client.impl.AccessToken(requestToken, requestSecret, verifier)
}

// AuthorizationURL blabla
func (client *oauth1Client) AuthorizationURL(requestToken string) (*url.URL, error) {
	return client.impl.AuthorizationURL(requestToken)
}

// HTTPClient blabla
func (client *oauth1Client) HTTPClient(accessToken, accessSecret string) *http.Client {
	config := client.impl
	token := oauth1.NewToken(accessToken, accessSecret)
	return config.Client(oauth1.NoContext, token)
}

// RequestToken blabla
func (client *oauth1Client) RequestToken() (requestToken, requestSecret string, err error) {
	return client.impl.RequestToken()
}

// ParseAuthorizationCallback blabla
func (client *oauth1Client) ParseAuthorizationCallback(r *http.Request) (requestToken, verifier string, err error) {
	return oauth1.ParseAuthorizationCallback(r)
}
