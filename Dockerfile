FROM golang:1.10.3
LABEL Author="Abdullah Alansari <ahimta@gmail.com>"
LABEL Name="tweeters-stats-golang"

RUN go get -u github.com/golang/dep/cmd/dep && \
  go get -u github.com/golang/lint/golint && \
  go get github.com/pilu/fresh

COPY Gopkg.lock Gopkg.toml /go/src/github.com/Ahimta/tweeters-stats-golang/
WORKDIR /go/src/github.com/Ahimta/tweeters-stats-golang

COPY . /go/src/github.com/Ahimta/tweeters-stats-golang
RUN dep ensure
RUN go build main.go

ENV CONSUMER_KEY consumerKey
ENV CONSUMER_SECRET consumerSecret
ENV CALLBACK_URL https://tweeters-stats.herokuapp.com
ENV PORT 8080
ENV HOMEPAGE /
ENV HOST tweeters-stats.herokuapp.com
ENV PROTOCOL https

EXPOSE 8080
RUN ./test
CMD ./main