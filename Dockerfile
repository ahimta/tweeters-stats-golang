FROM golang:1.9.2
LABEL Author="Abdullah Alansari <ahimta@gmail.com>"
LABEL Name="tweeters-stats-golang"

COPY . /go/src/github.com/Ahimta/tweeters-stats-golang
WORKDIR /go/src/github.com/Ahimta/tweeters-stats-golang

RUN go get -u github.com/golang/dep/cmd/dep
RUN go get github.com/pilu/fresh
RUN dep ensure
RUN go build

EXPOSE 8080

CMD fresh