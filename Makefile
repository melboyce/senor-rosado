SHELL := /bin/bash

TARGET := $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL: $(TARGET)

VERSION := 0.1.0
BUILD := `git rev-parse HEAD`

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean install uninstall fmt simplify check run

all: check test $(TARGET)

$(TARGET):
	go get -d ./...
	go build $(LDFLAGS) -o $(TARGET)

build: $(TARGET)
	@true

clean:
	rm -f $(TARGET)
	go clean

install:
	go install $(LDFLAGS)

uninstall: clean
	rm -f $$(which ${TARGET})

fmt:
	gofmt -l -w $(SRC)

simplify:
	gofmt -s -l -w $(SRC)

check:
	test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	go tool vet ${SRC}

run: install
	$(TARGET)

test:
	go get -d ./...
	go test ./...

static:
	go get -d ./...
	CGO_ENABLED=0 go build $(LDFLAGS) -o senor-rosado-static
