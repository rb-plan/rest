#!/bin/sh

mkdir -p target

# Set Go modules
export GO111MODULE=on

# Build for arm64-linux
GOARCH=arm64 GOOS=linux go build -o target/rest-arm64-linux main.go host_linux_arm64.go

# Build for arm64-freebsd
GOARCH=arm64 GOOS=freebsd go build -o target/rest-arm64-freebsd main.go host_linux_arm64.go

# Build for x86_64-linux
GOARCH=amd64 GOOS=linux go build -o target/rest-x86_64-linux main.go host_linux_arm64.go

# Build for x86_64-freebsd
GOARCH=amd64 GOOS=freebsd go build -o target/rest-x86_64-freebsd main.go host_linux_arm64.go

echo "Build complete."
