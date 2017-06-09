SHELL := /bin/bash
.DEFAULT_GOAL := help

.PHONY: carts run help

carts:
	@go build -buildmode=plugin -o plugins/hello.so _cartridges/hello.go
	@go build -buildmode=plugin -o plugins/weather.so _cartridges/weather.go

run:
	@go build && ./senor-rosado

help:
	@echo "make {carts,run}"
