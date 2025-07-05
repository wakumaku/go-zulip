#!/bin/sh
set -e pipefail 

echo "Running go mod tidy..."
go mod tidy

echo "Running linters..."
go tool golangci-lint run

echo "Running go test..."
go test -v --race --cover ./...