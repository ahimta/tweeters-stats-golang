FROM golang:1.9.2
LABEL Author="Abdullah Alansari <ahimta@gmail.com>"
LABEL Name="tweeters-stats-golang"

ARG livereload

COPY . /go/src/github.com/Ahimta/tweeters-stats-golang
WORKDIR /go/src/github.com/Ahimta/tweeters-stats-golang

RUN go get -u github.com/golang/dep/cmd/dep && \
  go get -u github.com/golang/lint/golint

RUN if [ ${livereload} = enabled ]; then go get github.com/pilu/fresh; fi

RUN dep ensure
RUN go build main.go

EXPOSE 8080

# @hack: can't use build argument in CMD so an environment variable is added
ENV LIVERELOAD ${livereload}
CMD if [ ${LIVERELOAD} = enabled ]; then fresh; else ./main; fi