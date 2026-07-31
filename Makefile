build:
	go build -o bin/gen-diff ./cmd/gen-diff

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

test:
	go test -v -race ./...

.PHONY: build lint lint-fix test 