# Web build image
FROM node:lts AS webbuilder

WORKDIR /usr/src/stage

COPY package.json .
COPY package-lock.json .
COPY tsconfig.json .
COPY tslint.json .

RUN npm install --only=prod --quiet

COPY src ./src
COPY public ./public

RUN npm run build

# Server build image
FROM golang:alpine AS gobuilder

# Install ca-certificates
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
# Create appuser.
ENV USER=appuser
ENV UID=10001 
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
    
WORKDIR $GOPATH/src/github.com/kubelens/kubelens/web

COPY ./public/server.go .

RUN go mod init && go mod tidy

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /go/bin/kubelens-web

# Run image
FROM scratch

# Import certs
COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the user and group files from the builder.
COPY --from=gobuilder /etc/passwd /etc/passwd
COPY --from=gobuilder /etc/group /etc/group

# Copy static executable.
COPY --from=gobuilder /go/bin/kubelens-web /go/bin/kubelens-web
# Copy static web files.
COPY --from=webbuilder /usr/src/stage/build /usr/src/website

# Use an unprivileged user.
USER appuser:appuser

# Run the binary.
ENTRYPOINT ["/go/bin/kubelens-web", "-confpath=/mnt/config", "-sitepath=/usr/src/website"]
