package services

import (
	"net/http"

	"github.com/Ahimta/tweeters-stats-golang/auth"
	"github.com/Ahimta/tweeters-stats-golang/entities"
	"github.com/dghubble/go-twitter/twitter"
)

// TweetsService blablabla
type TweetsService interface {
	FetchTweeters(accessToken, accessSecret string) ([]*entities.Tweeter, error)
}

type tweetsService struct {
	httpClientImpl func(accessToken, accessSecret string) (*http.Client, error)
}

// NewTweetsService blablabla
func NewTweetsService(oauthClient auth.Oauth1Client) TweetsService {
	return &tweetsService{oauthClient.HTTPClient}
}

// FetchTweeters blablabla
func (_tweetsService *tweetsService) FetchTweeters(accessToken, accessSecret string) ([]*entities.Tweeter, error) {
	httpClient, err := _tweetsService.httpClientImpl(accessToken, accessSecret)

	if err != nil {
		return nil, err
	}

	twitterClient := twitter.NewClient(httpClient)
	tweets, _, err := twitterClient.Timelines.HomeTimeline(&twitter.HomeTimelineParams{Count: 200})

	if err != nil {
		return nil, err
	}

	tweeters := make([]*entities.Tweeter, 0, len(tweets))
	for _, tweeter := range tweets {
		tweeters = append(tweeters, &entities.Tweeter{FullName: tweeter.User.Name, Username: tweeter.User.ScreenName})
	}

	return tweeters, nil
}
