# Tweeters Stats Golang
[![Build Status](https://travis-ci.org/Ahimta/tweeters-stats-golang.svg?branch=master)](https://travis-ci.org/Ahimta/tweeters-stats-golang)
[![Go Report Card](https://goreportcard.com/badge/Ahimta/tweeters-stats-golang)](https://goreportcard.com/report/Ahimta/tweeters-stats-golang)

## Requirements
* Twitter App (with read-only and login privileges)
* Docker
* Terraform (for deployment)

## Environment Variables
* CONSUMER_KEY: Twitter's consumer key
* CONSUMER_SECRET: Twitter's consumer secret
* CALLBACK_URL: Twitter's callback URL
* PORT: Port to use for the web server
* HOMEPAGE: URL to redirect to (e.g., when Twitter login successful)
* HOST: mostly for CSRF middleware
* PROTOCOL: mostly for CSRF middleware
* CORS_DOMAIN?: Domain to allow CORS (can useful for development)

## Test
* `docker build --file Dockerfile.test --tag tweeters-stats-golang-test .`
* `docker run -it --rm --env-file .env --env NEW_RELIC_LICENSE_KEY= -v $PWD:/go/src/github.com/Ahimta/tweeters-stats-golang tweeters-stats-golang-test`

## Build & Run (development)
* `docker build --file Dockerfile.dev --tag tweeters-stats-golang-dev .`
* `docker run -it --rm --env-file .env --env NEW_RELIC_LICENSE_KEY= -p 8080:8080 -v $PWD:/go/src/github.com/Ahimta/tweeters-stats-golang tweeters-stats-golang-dev`

## Build & Run (production)
* `docker build --file Dockerfile.prod --tag tweeters-stats-golang-prod .`
* `docker run -it --rm --env-file .env -p 8080:8080 tweeters-stats-golang-prod`

## Infrastructure (AWS)
* Make sure you have a default AWS profile configured
* Create `terraform.tfvars` file with Terraform variables
* `terraform init`
* `terraform apply`

## Deploy
* `aws ecr get-login --no-include-email --region eu-west-1 | bash`
* `docker build --file Dockerfile.prod --tag tweeters-stats-golang-prod .`
* `docker tag tweeters-stats-golang-prod:latest <ecr-repo>:latest`
* `docker push <ecr-repo>:latest`
* `aws ecs update-service --cluster backend --service backend --force-new-deployment`

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