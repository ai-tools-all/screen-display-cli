.PHONY: build clean install test

BINARY_NAME=dmon
BUILD_DIR=.
INSTALL_DIR=$(HOME)/.local/bin

build:
	go build -o $(BINARY_NAME) .

clean:
	rm -f $(BINARY_NAME) dmon-cli

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/

test:
	go test ./...

run: build
	./$(BINARY_NAME)

.DEFAULT_GOAL := build
