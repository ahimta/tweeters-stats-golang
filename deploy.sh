#!/usr/bin/env sh
heroku login
heroku container:login

docker build --file Dockerfile --tag tweeters-stats-golang .
docker tag tweeters-stats-golang registry.heroku.com/tweeters-stats/web

docker push registry.heroku.com/tweeters-stats/web
heroku container:release --app tweeters-stats web
