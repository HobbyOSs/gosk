# Go parameters
GOBUILD=go build
BIN=gosk
NASK=wine nask.exe
KILL_DEAD_CODE = find . -type f -name "*.go" -exec sed -i -E '/^\s*\/\/.*(remove|delete|unnecessary|dead code|no longer needed)/Id' {} +

.PHONY: all test gen

all: build test

build: gen
	go build ./...
	$(GOBUILD) -v -o $(BIN) ./cmd/gosk

test: dep
	go tool gotest -v ./...

clean:
	go clean

run: build
	./$(BIN)

gen: dep
	go generate ./...

fmt:
	$(KILL_DEAD_CODE)
	go fmt ./...

dep:
	go mod download
	go mod tidy
