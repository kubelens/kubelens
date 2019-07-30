# build image
FROM golang:latest AS build-env

WORKDIR /go/src/github.com/kubelens/kubelens/web

COPY ./build ./build
COPY ./public/server.go .

RUN rm -rf ./build/config ./build/server*

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -a -o ./server-x86-x64 server.go

# run image
FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk add ca-certificates
RUN update-ca-certificates

RUN addgroup -g 1000 -S dockeruser \
  && adduser -u 1000 -S dockeruser -G dockeruser

RUN mkdir -p /srv/site

COPY --from=build-env /go/src/github.com/kubelens/kubelens/web/ /tmp

RUN mv -v /tmp/server-x86-x64 /srv/
RUN mv -v /tmp/build/* /srv/site
RUN rm -rf /tmp/*

RUN chown -R dockeruser:dockeruser /srv/site

USER dockeruser

# /mnt/config must equal the mount path of the config map
ENTRYPOINT ["/srv/server-x86-x64", "-confpath=/mnt/config", "-sitepath=/srv/site"]
