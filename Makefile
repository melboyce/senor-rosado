SHELL := /bin/zsh

TARGET := $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL: $(TARGET)

VERSION := 0.1.0
BUILD := `git rev-parse HEAD`

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean install uninstall fmt simplify check run

all: check $(TARGET)

$(TARGET):
	go build $(LDFLAGS) -o $(TARGET)
	find _cartridges -name '*.go' | while read -r; do \
	go build -buildmode=plugin \
	-o plugins/$${$${REPLY##*/}%go}so $$REPLY; done

build: $(TARGET)
	@true

clean:
	rm -f $(TARGET)
	rm -f plugins/*
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
	for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done
	go tool vet ${SRC}

run: install
	$(TARGET)
