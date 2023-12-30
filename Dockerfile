# Use golang alpine flavor as our builder image
FROM golang:alpine3.18 as builder
# Add maintainer metadata
LABEL maintainer="aeolyus"
# Set our workspace
WORKDIR /app
# Required to compile C code (go-sqlite3)
RUN apk add --no-cache git gcc musl-dev
# Copy go.mod and go.sum dependency list into workspace
COPY go.mod go.sum ./
# Download dependencies; cached as long as we don't change mod/sum
RUN go mod download
# Copy over source code
COPY . .
# Go creates static binaries by default unless C code is called (go-sqlite3) in
# which case, it will create a dynamically linked binary. These flags will
# explicitly build a static binary. We also use the dedicated RUN cache to
# cache build outputs.
RUN --mount=type=cache,target=/root/.cache/go-build \
  go build -ldflags '-w -extldflags=-static' -o gull

# Use busybox as a lightweight base image with minimal set of utilities for
# debugging and troubleshooting
FROM busybox
# Set our workspace
WORKDIR /
# Copy the external assets from our builder image
COPY --from=builder /app/public /public
# Copy the binary from our builder image
COPY --from=builder /app/gull .
# Expose data volume
VOLUME /data
# Expose service port
EXPOSE 8081
# Run the binary
CMD ["./gull"]
