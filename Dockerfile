# Development Container
FROM golang:1.11 AS builder

# Add the script
ADD . /go/src/populate
WORKDIR /go/src/populate

# Set importante variables
ENV CGO_ENABLED 0
ENV GOOS "linux"
ENV GOARCH=amd64

# Install dependencies and build the app
RUN go get populate
RUN go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/populate

# Production Container
FROM scratch

# Copy app from dev container.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/populate /

ENTRYPOINT ["/populate"]