# Go parameters
GOBUILD=go build
BIN=gosk
NASK=wine nask.exe
KILL_DEAD_CODE = find . -type f -name "*.go" -exec sed -i -E '/^\s*\/\/.*(remove|delete|unnecessary|dead code|no longer needed)/Id' {} +

.PHONY: all test gen
# Use go run to execute gotestsum, ensuring it uses the version from go.mod
GOTESTSUM = go run gotest.tools/gotestsum

.PHONY: all test test-rerun-fails test-ci gen clean run dep

all: build test

build: gen
	go build ./...
	$(GOBUILD) -v -o $(BIN) ./cmd/gosk

test: dep
	$(GOTESTSUM) --format short-verbose -- ./...

test-rerun-fails: dep
	$(GOTESTSUM) --rerun-fails --format short-verbose --packages="./..."

test-ci: dep
	$(GOTESTSUM) --junitfile report.xml --format dots -- ./...

clean:
	go clean
	rm -f report.xml .gotestsum.json

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
