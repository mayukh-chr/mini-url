# Build stage
FROM golang:1.21.7-alpine3.19 AS builder

# Install git and ca-certificates (needed for downloading modules)
RUN apk update && apk upgrade && apk add --no-cache git ca-certificates && update-ca-certificates

# Create appuser
RUN adduser -D -g '' appuser

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s -extldflags "-static"' -a -installsuffix cgo -o main .

# Final stage
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /build/main /go/bin/main

# Copy templates if they exist
COPY --from=builder /build/templates /templates

# Use an unprivileged user
USER appuser

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/go/bin/main"]
