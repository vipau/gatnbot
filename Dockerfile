# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.25 AS builder

ENV CGO_ENABLED=0

WORKDIR /build

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /gatnbot-bin

# Runtime stage - use minimal image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /gatnbot

# Copy only the binary from builder
COPY --from=builder /gatnbot-bin /gatnbot-bin

# Copy any required runtime files (adjust as needed)
# COPY glados/ /gatnbot/glados/
# COPY jvazquez/ /gatnbot/jvazquez/

CMD ["/gatnbot-bin"]
