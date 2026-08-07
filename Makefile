build:
	go build -o bin/gendiff ./cmd/gendiff

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

test:
	go test -v -race ./...

coverage:
	go test -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

coverage-html:
	go test -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: build lint lint-fix test coverage coverage-html
