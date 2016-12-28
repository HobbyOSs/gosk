# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

BIN=gosk
NASK=wine nask.exe

.PHONY: all testdata

all: dep build test

build: dep
	cd cmd/gosk && $(GOBUILD) -v
	cd cmd/f12copy && $(GOBUILD) -v
	cd cmd/f12format && $(GOBUILD) -v
	cd ..
	$(GOINSTALL) -v ./...

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)

run: build
	./$(BIN)

fmt:
	for go_file in `find . -name \*.go`; do \
		go fmt $${go_file}; \
	done

dep:
	$(GOMOD) download
	$(GOMOD) tidy

tool:
	$(GOINSTALL) -v github.com/rogpeppe/godef@latest
	$(GOINSTALL) -v github.com/nsf/gocode@latest
	$(GOINSTALL) -v golang.org/x/lint/golint@latest
	$(GOINSTALL) -v github.com/kisielk/errcheck@latest
	$(GOINSTALL) -v github.com/go-delve/delve/cmd/dlv@latest

testdata:
	$(NASK) testdata/byte-opcode.nas testdata/byte-opcode.obj testdata/byte-opcode.list
