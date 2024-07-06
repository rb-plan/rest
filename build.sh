#!/bin/sh

mkdir -p target

# Set Go modules
export GO111MODULE=on

# Build for arm64-linux
GOARCH=arm64 GOOS=linux go build -o target/rest-arm64-linux main.go

# Build for arm-linux
GOARCH=arm GOOS=linux go build -o target/rest-arm-linux main.go

# Build for arm64-freebsd
GOARCH=arm64 GOOS=freebsd go build -o target/rest-arm64-freebsd main.go

# Build for x86_64-linux
GOARCH=amd64 GOOS=linux go build -o target/rest-x86_64-linux main.go

# Build for x86_64-freebsd
GOARCH=amd64 GOOS=freebsd go build -o target/rest-x86_64-freebsd main.go

echo "Build complete."
