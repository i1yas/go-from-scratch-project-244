build:
	go build -o bin/gendiff ./cmd/gendiff

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

test:
	go test -race -coverpkg=./...  -coverprofile=coverage.out ./...

coverage: test
	go tool cover -func=coverage.out

coverage-html: test
	go tool cover -html=coverage.out

.PHONY: build lint lint-fix test coverage coverage-html
