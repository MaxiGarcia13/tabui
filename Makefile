.PHONY: install lint test

install:
	go install ./cmd/tabui

lint:
	golangci-lint run ./...

test:
	go test ./...
