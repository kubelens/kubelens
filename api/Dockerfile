# Build image
FROM golang:alpine AS builder

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
    
WORKDIR $GOPATH/src/github.com/kubelens/kubelens/api
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /go/bin/kubelens-api

# Run image
FROM scratch

# Import certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable.
COPY --from=builder /go/bin/kubelens-api /go/bin/kubelens-api

# Use an unprivileged user.
USER appuser:appuser

# Run the binary.
ENTRYPOINT ["/go/bin/kubelens-api"]