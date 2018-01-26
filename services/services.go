package services

import (
	"errors"
	"net/http"

	"github.com/Ahimta/tweeters-stats-golang/auth"
	"github.com/Ahimta/tweeters-stats-golang/entities"
	"github.com/dghubble/go-twitter/twitter"
)

// TweetsService blablabla
type TweetsService interface {
	Tweeters(accessToken, accessSecret string) ([]*entities.Tweeter, error)
}

type tweetsService struct {
	tweetsImpl     func(httpClient *http.Client) ([]twitter.Tweet, error)
	httpClientImpl func(accessToken, accessSecret string) (*http.Client, error)
}

// NewTweetsService blablabla
func NewTweetsService(client auth.Oauth1Client) TweetsService {
	return &tweetsService{getTweets, client.HTTPClient}
}

// Tweeters blablabla
func (service *tweetsService) Tweeters(
	accessToken,
	accessSecret string,
) ([]*entities.Tweeter, error,
) {

	if accessToken == "" || accessSecret == "" {
		return nil, errors.New("services: missing accessToken or accessSecret")
	}

	httpClient, err := service.httpClientImpl(accessToken, accessSecret)

	if err != nil {
		return nil, err
	}

	tweets, err := service.tweetsImpl(httpClient)

	if err != nil {
		return nil, err
	}

	tweeters := make([]*entities.Tweeter, 0, len(tweets))
	for _, tweeter := range tweets {
		tweeters = append(
			tweeters,
			&entities.Tweeter{
				FullName: tweeter.User.Name,
				Username: tweeter.User.ScreenName,
			})
	}

	return tweeters, nil
}

func getTweets(client *http.Client) ([]twitter.Tweet, error) {
	twitterClient := twitter.NewClient(client)
	tweets, _, err := twitterClient.
		Timelines.
		HomeTimeline(&twitter.HomeTimelineParams{Count: 200})

	if err != nil {
		return nil, err
	}

	return tweets, nil
}
