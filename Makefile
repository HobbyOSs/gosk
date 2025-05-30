# Go parameters
GOBUILD=go build
GOTEST=gotest

BIN=gosk
NASK=wine nask.exe
KILL_DEAD_CODE = find . -type f -name "*.go" -exec sed -i -E '/^\s*\/\/.*(remove|delete|unnecessary|dead code|no longer needed)/Id' {} +

.PHONY: all test gen

all: build test

build: dep gen
	go build ./...
	$(GOBUILD) -v -o $(BIN) ./cmd/gosk

test:
	go install -v github.com/rakyll/gotest@latest
	$(GOTEST) -v ./...

clean:
	go clean

run: build
	./$(BIN)

gen:
	go generate ./...

fmt:
	$(KILL_DEAD_CODE)
	go fmt ./...

dep:
	go install -v github.com/mna/pigeon@latest
	go install -v github.com/Bin-Huang/newc@latest
	go install -v github.com/dmarkham/enumer@latest
	go mod download
	go mod tidy

tool:
	go install -v golang.org/x/tools/gopls@latest
	go install -v github.com/go-delve/delve/cmd/dlv@latest
	go install -v github.com/awalterschulze/goderive@latest
	go install -v github.com/hairyhenderson/gomplate/v3/cmd/gomplate@latest
