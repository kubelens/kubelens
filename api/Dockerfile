# build image
FROM golang:latest AS build-env
WORKDIR /go/src/github.com/kubelens/kubelens/api/ 
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -a -o ./bin/server-x86-x64 *.go

# run image
FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk add ca-certificates
RUN update-ca-certificates

RUN addgroup -g 1000 -S dockeruser \
  && adduser -u 1000 -S dockeruser -G dockeruser

RUN mkdir -p /srv

COPY --from=build-env /go/src/github.com/kubelens/kubelens/api/bin/ /srv

USER dockeruser

ENTRYPOINT ["/srv/server-x86-x64"]
