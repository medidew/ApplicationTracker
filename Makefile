BINARY := bin/ApplicationTracker
GO := go

.DEFAULT_GOAL := all

.PHONY: all test build clean fmt vet

all: test build

test:
	$(GO) test ./... -v

build:
	mkdir -p $(dir $(BINARY))
	$(GO) build -o $(BINARY) .

clean:
	rm -rf bin

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...
