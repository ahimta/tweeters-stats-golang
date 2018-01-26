package usecases

import (
	"errors"
	"net/http"
	"net/url"
	"sort"

	"github.com/Ahimta/tweeters-stats-golang/auth"
	"github.com/Ahimta/tweeters-stats-golang/entities"
	"github.com/Ahimta/tweeters-stats-golang/services"
)

// LoginResult blablabla
type LoginResult struct {
	AuthorizationURL *url.URL
	RequestSecret    string
}

// Oauth1CallbackResult blablabla
type Oauth1CallbackResult struct {
	AccessToken  string
	AccessSecret string
}

type statsSort []*entities.TweeterStats

func (xs statsSort) Len() int {
	return len(xs)
}

func (xs statsSort) Swap(i, j int) {
	xs[i], xs[j] = xs[j], xs[i]
}

func (xs statsSort) Less(i, j int) bool {
	return (xs[i].TweetsCount < xs[j].TweetsCount)
}

// TweetersStats blablabla
func TweetersStats(
	tweetsService services.TweetsService,
	accessToken,
	accessSecret string,
) (
	[]*entities.TweeterStats, error,
) {

	if accessToken == "" || accessSecret == "" {
		return nil, errors.New("usecases: accessToken or accessSecret missing -_-")
	}

	tweeters, err := tweetsService.Tweeters(accessToken, accessSecret)

	if err != nil {
		return nil, err
	}

	statsByUsername := make(map[string]*entities.TweeterStats)
	for _, tweeter := range tweeters {
		tweeterStats, ok := statsByUsername[tweeter.Username]

		if ok {
			tweeterStats.TweetsCount++
		} else {
			statsByUsername[tweeter.Username] = &entities.TweeterStats{
				FullName:    tweeter.FullName,
				Username:    tweeter.Username,
				TweetsCount: 1,
			}
		}
	}

	tweetersStats := make([]*entities.TweeterStats, 0, len(statsByUsername))

	for _, tweeterStats := range statsByUsername {
		tweetersStats = append(tweetersStats, &entities.TweeterStats{
			FullName:    tweeterStats.FullName,
			Username:    tweeterStats.Username,
			TweetsCount: tweeterStats.TweetsCount,
		})
	}

	sort.Sort(sort.Reverse(statsSort(tweetersStats)))
	return tweetersStats, nil
}

// Oauth1Callback blablabla
func Oauth1Callback(
	client auth.Oauth1Client,
	requestSecret string,
	r *http.Request,
) (
	*Oauth1CallbackResult, error,
) {

	if requestSecret == "" || r == nil {
		return nil, errors.New("usecases: requestSecret or request missing -_-")
	}

	requestToken, verifier, err := client.ParseAuthorizationCallback(r)

	if err != nil {
		return nil, err
	}

	accessToken, accessSecret, err := client.AccessToken(
		requestToken,
		requestSecret,
		verifier,
	)

	if err != nil {
		return nil, err
	}

	return &Oauth1CallbackResult{
		AccessToken:  accessToken,
		AccessSecret: accessSecret,
	}, nil
}

// Login blablabla
func Login(client auth.Oauth1Client) (*LoginResult, error) {
	requestToken, requestSecret, err := client.RequestToken()

	if err != nil {
		return nil, err
	}

	authorizationURL, err := client.AuthorizationURL(requestToken)
	return &LoginResult{
		AuthorizationURL: authorizationURL,
		RequestSecret:    requestSecret,
	}, nil
}
