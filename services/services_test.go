package services

import (
	"reflect"
	"testing"

	"github.com/Ahimta/tweeters-stats-golang/auth"
	"github.com/Ahimta/tweeters-stats-golang/entities"
)

func TestNewTweetsService(t *testing.T) {
	oauth1Client, _ := auth.NewOauth1Client("consumerKey", "consumerSecret", "callbackURL")

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
			want: &tweetsService{oauth1Client.HTTPClient},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTweetsService(tt.args.oauthClient)
			gotHTTPClientImplPointer := reflect.ValueOf(got).Elem().FieldByName("httpClientImpl").Pointer()
			wantHTTPClientImplPointer := reflect.ValueOf(tt.want).Elem().FieldByName("httpClientImpl").Pointer()

			if gotHTTPClientImplPointer != wantHTTPClientImplPointer {
				t.Errorf("NewTweetsService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tweetsService_FetchTweeters(t *testing.T) {
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
	// TODO: Add test cases.
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
