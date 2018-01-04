package auth

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/dghubble/oauth1"
)

var validClient, _ = NewOauth1Client("consumerKey", "consumerSecret", "callbackURL")

func TestNewOauth1Client(t *testing.T) {
	//
	client, err := NewOauth1Client("blablabla", "blablabla", "blablabla")

	if err != nil {
		t.Errorf(err.Error())
	}

	//
	client, err = NewOauth1Client("blablabla", "", "blablabla")

	if err == nil || client != nil {
		t.Errorf("should return an error when a required config value is missing!")
	}
}

func Test_oauth1Client_AccessToken(t *testing.T) {
	type args struct {
		requestToken  string
		requestSecret string
		verifier      string
	}
	tests := []struct {
		name             string
		client           *oauth1Client
		args             args
		wantAccessToken  string
		wantAccessSecret string
		wantErr          bool
	}{
		{
			name:    "should return an error when requestToken, requestSecret, or verifier is missing",
			client:  &oauth1Client{},
			args:    args{"blablabla", "", "blablabla"},
			wantErr: true,
		},
		{
			name: "should call the actual implementation with arguments",
			client: &oauth1Client{
				accessTokenImpl: func(requestToken, requestSecret, verifier string) (accessToken, accessSecret string, err error) {
					if requestToken != "requestToken" || requestSecret != "requestSecret" || verifier != "verifier" {
						t.Errorf("expected to call actual implementation with args -_-")
					}

					return "accessTokenResult", "accessSecretResult", nil
				},
			},
			args: args{"requestToken", "requestSecret", "verifier"},

			wantAccessToken:  "accessTokenResult",
			wantAccessSecret: "accessSecretResult",
		},
		{
			name: "should return error from actual implementation",
			client: &oauth1Client{
				accessTokenImpl: func(requestToken, requestSecret, verifier string) (accessToken, accessSecret string, err error) {
					if requestToken != "requestToken" || requestSecret != "requestSecret" || verifier != "verifier" {
						t.Errorf("expected to call actual implementation with args -_-")
					}

					return "", "", errors.New("whaaat -_-")
				},
			},

			args:    args{"requestToken", "requestSecret", "verifier"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAccessToken, gotAccessSecret, err := tt.client.AccessToken(tt.args.requestToken, tt.args.requestSecret, tt.args.verifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("oauth1Client.AccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAccessToken != tt.wantAccessToken {
				t.Errorf("oauth1Client.AccessToken() gotAccessToken = %v, want %v", gotAccessToken, tt.wantAccessToken)
			}
			if gotAccessSecret != tt.wantAccessSecret {
				t.Errorf("oauth1Client.AccessToken() gotAccessSecret = %v, want %v", gotAccessSecret, tt.wantAccessSecret)
			}
		})
	}
}

func Test_oauth1Client_AuthorizationURL(t *testing.T) {
	urlResult := &url.URL{}

	type args struct {
		requestToken string
	}
	tests := []struct {
		name    string
		client  *oauth1Client
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			name: "should return actual implementation result",
			client: &oauth1Client{
				authorizationURLImpl: func(requestToken string) (*url.URL, error) {
					return urlResult, nil
				},
			},

			args: args{"requestToken"},
			want: urlResult,
		},
		{
			name: "should pass and return actual implementation error",
			client: &oauth1Client{
				authorizationURLImpl: func(requestToken string) (*url.URL, error) {
					return nil, errors.New("whaaat -_-")
				},
			},

			args:    args{"requestToken"},
			wantErr: true,
		},
		{
			name:    "should return an error when requestToken is missing",
			args:    args{""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.AuthorizationURL(tt.args.requestToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("oauth1Client.AuthorizationURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("oauth1Client.AuthorizationURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_oauth1Client_HTTPClient(t *testing.T) {
	client := &http.Client{}
	token := &oauth1.Token{}

	type args struct {
		accessToken  string
		accessSecret string
	}
	tests := []struct {
		name    string
		client  *oauth1Client
		args    args
		want    *http.Client
		wantErr bool
	}{
		{
			name: "should pass values to actual implementation",
			client: &oauth1Client{
				newTokenImpl: func(accessToken, accessSecret string) *oauth1.Token {
					if accessToken != "accessToken" || accessSecret != "accessSecret" {
						t.Errorf("Whaaat!")
					}

					return token
				},
				clientImpl: func(ctx context.Context, t0 *oauth1.Token) *http.Client {
					if ctx != oauth1.NoContext || t0 != token {
						t.Errorf("Whaaat!")
					}

					return client
				},
			},
			args: args{"accessToken", "accessSecret"},
			want: client,
		},
		{
			name:    "should return an error when accessToken or accessSecret is missing",
			args:    args{"accesToken", ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.HTTPClient(tt.args.accessToken, tt.args.accessSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("oauth1Client.HTTPClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("oauth1Client.HTTPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_oauth1Client_RequestToken(t *testing.T) {
	tests := []struct {
		name              string
		client            *oauth1Client
		wantRequestToken  string
		wantRequestSecret string
		wantErr           bool
	}{
		{
			name: "should pass and return values using actual implementation",
			client: &oauth1Client{
				requestTokenImpl: func() (requestToken, requestSecret string, err error) {
					return "requestToken", "requestSecret", nil
				},
			},

			wantRequestToken:  "requestToken",
			wantRequestSecret: "requestSecret",
		},
		{
			name: "should return actual implementation error",
			client: &oauth1Client{
				requestTokenImpl: func() (requestToken, requestSecret string, err error) {
					return "", "", errors.New("whaaat -_-")
				},
			},

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRequestToken, gotRequestSecret, err := tt.client.RequestToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("oauth1Client.RequestToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRequestToken != tt.wantRequestToken {
				t.Errorf("oauth1Client.RequestToken() gotRequestToken = %v, want %v", gotRequestToken, tt.wantRequestToken)
			}
			if gotRequestSecret != tt.wantRequestSecret {
				t.Errorf("oauth1Client.RequestToken() gotRequestSecret = %v, want %v", gotRequestSecret, tt.wantRequestSecret)
			}
		})
	}
}

func Test_oauth1Client_ParseAuthorizationCallback(t *testing.T) {
	request := &http.Request{}

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name             string
		client           *oauth1Client
		args             args
		wantRequestToken string
		wantVerifier     string
		wantErr          bool
	}{
		{
			name: "should pass and return values using actual implementation",
			client: &oauth1Client{
				parseAuthorizationCallbackImpl: func(r *http.Request) (requestToken, verifier string, err error) {
					if r != request {
						t.Errorf("Whaaat!")
					}

					return "requestToken", "verifier", nil
				},
			},

			args:             args{request},
			wantRequestToken: "requestToken",
			wantVerifier:     "verifier",
		},
		{
			name: "should return error from actual implementation when appropriate",
			client: &oauth1Client{
				parseAuthorizationCallbackImpl: func(r *http.Request) (requestToken, verifier string, err error) {
					return "", "", errors.New("whaaat -_-")
				},
			},

			wantErr: true,
		},
		{
			name:    "should return an error when r is missing",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRequestToken, gotVerifier, err := tt.client.ParseAuthorizationCallback(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("oauth1Client.ParseAuthorizationCallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRequestToken != tt.wantRequestToken {
				t.Errorf("oauth1Client.ParseAuthorizationCallback() gotRequestToken = %v, want %v", gotRequestToken, tt.wantRequestToken)
			}
			if gotVerifier != tt.wantVerifier {
				t.Errorf("oauth1Client.ParseAuthorizationCallback() gotVerifier = %v, want %v", gotVerifier, tt.wantVerifier)
			}
		})
	}
}
