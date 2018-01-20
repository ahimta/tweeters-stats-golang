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

// OauthLoginResult blablabla
type OauthLoginResult struct {
	AuthorizationURL *url.URL
	RequestSecret    string
}

// HandleOauth1CallbackResult blablabla
type HandleOauth1CallbackResult struct {
	AccessToken  string
	AccessSecret string
}

type tweeterStatsSort []*entities.TweeterStats

func (xs tweeterStatsSort) Len() int {
	return len(xs)
}

func (xs tweeterStatsSort) Swap(i, j int) {
	xs[i], xs[j] = xs[j], xs[i]
}

func (xs tweeterStatsSort) Less(i, j int) bool {
	return (xs[i].TweetsCount < xs[j].TweetsCount)
}

// GetTweetersStats blablabla
func GetTweetersStats(
	tweetsService services.TweetsService,
	accessToken,
	accessSecret string,
) (
	[]*entities.TweeterStats, error,
) {

	if accessToken == "" || accessSecret == "" {
		return nil, errors.New("usecases: accessToken or accessSecret missing -_-")
	}

	tweeters, err := tweetsService.FetchTweeters(accessToken, accessSecret)

	if err != nil {
		return nil, err
	}

	tweetersStatsByUsername := make(map[string]*entities.TweeterStats)
	for _, tweeter := range tweeters {
		tweeterStats, ok := tweetersStatsByUsername[tweeter.Username]

		if ok {
			tweeterStats.TweetsCount++
		} else {
			tweetersStatsByUsername[tweeter.Username] = &entities.TweeterStats{
				FullName:    tweeter.FullName,
				Username:    tweeter.Username,
				TweetsCount: 1,
			}
		}
	}

	tweetersStats := make(
		[]*entities.TweeterStats,
		0,
		len(tweetersStatsByUsername),
	)

	for _, tweeterStats := range tweetersStatsByUsername {
		tweetersStats = append(tweetersStats, &entities.TweeterStats{
			FullName:    tweeterStats.FullName,
			Username:    tweeterStats.Username,
			TweetsCount: tweeterStats.TweetsCount,
		})
	}

	sort.Sort(sort.Reverse(tweeterStatsSort(tweetersStats)))
	return tweetersStats, nil
}

// HandleOauth1Callback blablabla
func HandleOauth1Callback(
	oauthClient auth.Oauth1Client,
	requestSecret string,
	r *http.Request,
) (
	*HandleOauth1CallbackResult, error,
) {

	if requestSecret == "" || r == nil {
		return nil, errors.New("usecases: requestSecret or request missing -_-")
	}

	requestToken, verifier, err := oauthClient.ParseAuthorizationCallback(r)

	if err != nil {
		return nil, err
	}

	accessToken, accessSecret, err := oauthClient.AccessToken(
		requestToken,
		requestSecret,
		verifier,
	)

	if err != nil {
		return nil, err
	}

	return &HandleOauth1CallbackResult{
		AccessToken:  accessToken,
		AccessSecret: accessSecret,
	}, nil
}

// Login blablabla
func Login(client auth.Oauth1Client) (*OauthLoginResult, error) {
	requestToken, requestSecret, err := client.RequestToken()

	if err != nil {
		return nil, err
	}

	authorizationURL, err := client.AuthorizationURL(requestToken)
	return &OauthLoginResult{
		AuthorizationURL: authorizationURL,
		RequestSecret:    requestSecret,
	}, nil
}
