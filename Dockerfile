# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.15 as builder

# Copy local code to the container image.
WORKDIR /go/src/github.com/hayashiki/tarsier-integration
COPY go.mod .
COPY go.sum .

# Get dependencies - will be cached if we don't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

RUN cd cmd/tarsier-integration && \
    CGO_ENABLED=0 GOOS=linux go build -a -v -o tarsier-integration

# Use a Docker multi-stage build to create a lean production image.
FROM alpine
RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/github.com/hayashiki/tarsier-integration/cmd/tarsier-integration/tarsier-integration /tarsier-integration

# Run the web service on container startup.
CMD ["/tarsier-integration"]
