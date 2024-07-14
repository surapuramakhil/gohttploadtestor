# Use the official Golang Alpine image to create a build artifact.
FROM golang:1.22.1-alpine as builder

# Install necessary packages
RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o loadtester .

# Run tests
FROM builder AS tester
RUN go test -v ./...

# Final stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/loadtester /app/loadtester
ENTRYPOINT ["/app/loadtester"]