echo "Running go mod tidy..."
go mod tidy

echo "Running linters..."
golangci-lint run

echo "Running go test..."
go test -v --race --cover ./...