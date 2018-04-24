default: build

build:
	@echo "Building exec-docker"
	@go build -ldflags="-s -w"

test: 
	@go test -coverprofile=coverage.out

view: 
	@go tool cover -html=coverage.out

.PHONY: build test view
