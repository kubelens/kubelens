# kubelens/api

## Setup

### Install Go

Official Docs: https://golang.org/doc/install

### Make

For Windows, you can choose to install Make if you want to run make commands without docker: https://chocolatey.org/packages/make

For Mac, you need XCode installed which includes make.

If you don't want to use make, you can run the commands from the Makefile directly as they just wrap script calls.

## Build & Deploy

### `DOCKER_ID=id DOCKER_USER=user GIT_BRANCH=master make docker-build-push`

This will build and and push the docker image.

### `make set-config`

Run this scripThe script assumes the required environment variables are set for the script (and so there doesn't have to be upteen args passed). Example (shortened): SERVER_PORT=39000 ALLOWED_ORIGINS='"http://kubelens.local","http://localhost:3000"' make set-config

### `INGRESS_HOST=api.kubelens.local make helm-upgrade`

Deploy via [Helm](https://helm.sh/)
