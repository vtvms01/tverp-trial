.PHONY: build vet test test-v cover tidy

build:
	go build ./...

vet:
	go vet ./...

# Run all tests
test:
	go test ./...

# Verbose test output
test-v:
	go test -v ./...

# Coverage report
cover:
	go test -cover ./...

tidy:
	go mod tidy
