# Tweeters Stats Golang
[![Build Status](https://travis-ci.org/Ahimta/tweeters-stats-golang.svg?branch=master)](https://travis-ci.org/Ahimta/tweeters-stats-golang)
[![Go Report Card](https://goreportcard.com/badge/Ahimta/tweeters-stats-golang)](https://goreportcard.com/report/Ahimta/tweeters-stats-golang)

## Requirements
* Twitter App (with read-only and login privileges)
* Docker

## Environment Variables
* CONSUMER_KEY: Twitter's consumer key
* CONSUMER_SECRET: Twitter's consumer secret
* CALLBACK_URL: Twitter's callback URL
* PORT: Port to use for the web server
* HOMEPAGE: URL to redirect to (e.g., when Twitter login successful)
* CORS_DOMAIN?: Domain to allow CORS (can useful for development)

## Build & Run (development)
* `docker build --tag tweeters-stats-golang --build-arg livereload=enabled .`
* `docker run -it --rm --env-file .env --env HOMEPAGE=/ -p 8080:8080 -v $PWD:/go/src/github.com/Ahimta/tweeters-stats-golang tweeters-stats-golang`

## Build & Run (production)
* `docker build --tag tweeters-stats-golang --build-arg livereload=disabled .`
* `docker run -it --rm --env-file .env --env HOMEPAGE=/ -p 8080:8080 -v $PWD:/go/src/github.com/Ahimta/tweeters-stats-golang tweeters-stats-golang`

## Test
* `docker run -it --rm --env-file .env --env HOMEPAGE=/ -p 8080:8080 -v $PWD:/go/src/github.com/Ahimta/tweeters-stats-golang tweeters-stats-golang ./test`

## Routes
* `/`: SPA frontend serving `index.html` (you have to provide your own)
* `/login/twitter`: Twitter's OAuth1 login
* `/oauth/twitter/callback`: Twitter's OAuth1 login callback
* `/tweeters-stats`: Tweeter's stats for authenticated Twitter account

## Recommended Development Environment
* OS: Ubuntu
* Editior: VS Code (using `Docker` and `Go` plugins)

## License
GNU General Public License v3.0