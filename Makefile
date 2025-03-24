.PHONY: build
build:
		go build -v ./cmd/apiServer

.DEFAULT_GOAL := build