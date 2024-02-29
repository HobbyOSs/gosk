# Go parameters
GOBUILD=go build
GOTEST=gotest

BIN=gosk
NASK=wine nask.exe

.PHONY: all testdata gen

all: dep build test

build: dep
	cd cmd/gosk && $(GOBUILD) -v
	cd cmd/f12copy && $(GOBUILD) -v
	cd cmd/f12format && $(GOBUILD) -v
	cd ..
	go install -v ./...

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
	go fmt ./...

dep:
	go mod download
	go mod tidy

tool:
	go install -v github.com/rogpeppe/godef@latest
	go install -v github.com/nsf/gocode@latest
	go install -v golang.org/x/lint/golint@latest
	go install -v github.com/kisielk/errcheck@latest
	go install -v github.com/go-delve/delve/cmd/dlv@latest
	go install -v golang.org/x/tools/gopls@latest
	go install -v github.com/mna/pigeon@latest

testdata:
	$(NASK) testdata/byte-opcode.nas testdata/byte-opcode.obj testdata/byte-opcode.list
