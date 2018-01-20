package services

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/Ahimta/tweeters-stats-golang/auth"
	"github.com/Ahimta/tweeters-stats-golang/entities"
	"github.com/dghubble/go-twitter/twitter"
)

func TestNewTweetsService(t *testing.T) {
	oauth1Client, _ := auth.NewOauth1Client(
		"consumerKey",
		"consumerSecret",
		"callbackURL",
	)

	type args struct {
		oauthClient auth.Oauth1Client
	}
	tests := []struct {
		name string
		args args
		want TweetsService
	}{
		{
			name: "should assign passed oauth1 client implementation",
			args: args{oauth1Client},
			want: &tweetsService{getTweets, oauth1Client.HTTPClient},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTweetsService(tt.args.oauthClient)
			gotHTTPClientImplPointer := reflect.
				ValueOf(got).
				Elem().
				FieldByName("httpClientImpl").
				Pointer()

			wantHTTPClientImplPointer := reflect.
				ValueOf(tt.want).
				Elem().
				FieldByName("httpClientImpl").
				Pointer()

			if gotHTTPClientImplPointer != wantHTTPClientImplPointer {
				t.Errorf("NewTweetsService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tweetsService_FetchTweeters(t *testing.T) {
	_httpClient := &http.Client{}

	type args struct {
		accessToken  string
		accessSecret string
	}
	tests := []struct {
		name           string
		_tweetsService *tweetsService
		args           args
		want           []*entities.Tweeter
		wantErr        bool
	}{
		{
			name: "should use the underlying implementation correctly",

			_tweetsService: &tweetsService{
				getTweetsImpl: func(httpClient *http.Client) ([]twitter.Tweet, error) {
					if httpClient != _httpClient {
						t.Errorf("httpClient not passed correctly")
					}

					return []twitter.Tweet{}, nil
				},
				httpClientImpl: func(accessToken, accessSecret string) (*http.Client, error) {
					if accessToken != "accessToken" || accessSecret != "accessSecret" {
						t.Errorf("accessToken or accessSecret not passed correctly")
					}

					return _httpClient, nil
				},
			},

			args: args{"accessToken", "accessSecret"},
			want: []*entities.Tweeter{},
		},
		{
			name: "should process the returned tweets correctly",

			_tweetsService: &tweetsService{
				getTweetsImpl: func(httpClient *http.Client) ([]twitter.Tweet, error) {
					return []twitter.Tweet{
						{
							User: &twitter.User{Name: "John Smith", ScreenName: "jsmith"},
						},
					}, nil
				},
				httpClientImpl: func(accessToken, accessSecret string) (*http.Client, error) {
					return _httpClient, nil
				},
			},

			args: args{"accessToken", "accessSecret"},
			want: []*entities.Tweeter{
				{
					FullName: "John Smith",
					Username: "jsmith",
				},
			},
		},
		{
			name: "should return an error when httpClientImpl does",

			_tweetsService: &tweetsService{
				getTweetsImpl: func(httpClient *http.Client) ([]twitter.Tweet, error) {
					return []twitter.Tweet{}, nil
				},
				httpClientImpl: func(accessToken, accessSecret string) (*http.Client, error) {
					return nil, errors.New("whaaat -_-")
				},
			},

			args:    args{"accessToken", "accessSecret"},
			wantErr: true,
		},
		{
			name: "should return an error when getTweetsImpl does",

			_tweetsService: &tweetsService{
				getTweetsImpl: func(httpClient *http.Client) ([]twitter.Tweet, error) {
					return nil, errors.New("whaaat -_-")
				},
				httpClientImpl: func(accessToken, accessSecret string) (*http.Client, error) {
					return _httpClient, nil
				},
			},

			args:    args{"accessToken", "accessSecret"},
			wantErr: true,
		},
		{
			name:           "should return an error when accessToken or accessSecret is missing",
			_tweetsService: &tweetsService{},
			args:           args{"accessToken", ""},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt._tweetsService.FetchTweeters(tt.args.accessToken, tt.args.accessSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("tweetsService.FetchTweeters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tweetsService.FetchTweeters() = %v, want %v", got, tt.want)
			}
		})
	}
}
