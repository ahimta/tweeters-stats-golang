# Tweeters Stats Golang

[![Go Report Card](https://goreportcard.com/badge/Ahimta/tweeters-stats-golang)](https://goreportcard.com/report/Ahimta/tweeters-stats-golang)
[![Build Status](https://travis-ci.org/Ahimta/tweeters-stats-golang.svg?branch=master)](https://travis-ci.org/Ahimta/tweeters-stats-golang)
[![Maintainability](https://api.codeclimate.com/v1/badges/9a3540991baf29bfc53b/maintainability)](https://codeclimate.com/github/Ahimta/tweeters-stats-golang/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/9a3540991baf29bfc53b/test_coverage)](https://codeclimate.com/github/Ahimta/tweeters-stats-golang/test_coverage)

## Requirements

- Twitter App (with only read-only permissions and no login privileges)
- Docker

## Environment Variables

- CONSUMER_KEY: Twitter's consumer key
- CONSUMER_SECRET: Twitter's consumer secret
- CALLBACK_URL: Twitter's callback URL
- PORT: Port to use for the web server
- HOMEPAGE: URL to redirect to (e.g., when Twitter login successful)
- HOST: mostly for CSRF middleware
- PROTOCOL: mostly for CSRF middleware
- CORS_DOMAIN?: Domain to allow CORS (can useful for development)

# Build

`docker build --file Dockerfile --tag tweeters-stats-golang .`

## Test

`docker run -it --rm --env-file .env --env NEW_RELIC_LICENSE_KEY= tweeters-stats-golang ./test`

## Run (local development)

1. `docker run -it --rm --env-file .env -v $PWD:/go/src/github.com/Ahimta/tweeters-stats-golang tweeters-stats-golang dep ensure`
2. `docker run -it --rm --env-file .env --env NEW_RELIC_LICENSE_KEY= -p 8080:8080 -v $PWD:/go/src/github.com/Ahimta/tweeters-stats-golang tweeters-stats-golang fresh`

## Run (production)

`docker run -it --rm --env-file .env -p 8080:8080 tweeters-stats-golang`

## Deploy

`sh deploy.sh`

## Routes

- `/`: SPA frontend serving `index.html` (you have to provide your own)
- `/login/twitter`: Twitter's OAuth1 login
- `/oauth/twitter/callback`: Twitter's OAuth1 login callback
- `/tweeters-stats`: Tweeter's stats for authenticated Twitter account

## Recommended Development Environment

- Editior: VS Code (using `Docker` and `Go` plugins)
- OS: Ubuntu

## License

GNU General Public License v3.0