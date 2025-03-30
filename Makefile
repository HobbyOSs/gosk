# Go parameters
GOBUILD=go build
GOTEST=gotest

BIN=gosk
NASK=wine nask.exe
KILL_DEAD_CODE = find . -type f -name "*.go" -exec sed -i -E '/^\s*\/\/.*(remove|delete|unnecessary|dead code|no longer needed)/Id' {} +

.PHONY: all test gen compress

all: build test

build: dep gen compress
	go build ./...
	cd cmd/gosk && $(GOBUILD) -v

test:
	go install -v github.com/rakyll/gotest@latest
	go test -cv ./...
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

compress:
	if [ ! -f pkg/asmdb/json-x86-64/x86_64.json.gz ] || [ pkg/asmdb/json-x86-64/x86_64.json -nt pkg/asmdb/json-x86-64/x86_64.json.gz ]; then \
		gzip -c pkg/asmdb/json-x86-64/x86_64.json > pkg/asmdb/json-x86-64/x86_64.json.gz; \
	fi

tool:
	go install -v golang.org/x/tools/gopls@latest
	go install -v github.com/go-delve/delve/cmd/dlv@latest
	go install -v github.com/awalterschulze/goderive@latest
	go install -v github.com/hairyhenderson/gomplate/v3/cmd/gomplate@latest
