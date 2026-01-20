.PHONY: build clean install test

BINARY_NAME=dmon
BUILD_DIR=.
INSTALL_DIR=$(HOME)/.local/bin
PKG_VERSION=github.com/abhishek/dmon-cli/internal/version

# Dynamic variables
GIT_COMMIT=$(shell git rev-parse --short HEAD 2> /dev/null || echo "none")
GIT_VERSION=$(shell git describe --tags --always --dirty 2> /dev/null || echo "dev")
BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Linker flags
# -s: strip symbol table
# -w: strip DWARF debug info
LDFLAGS=-ldflags "-s -w \
	-X '$(PKG_VERSION).GitCommit=$(GIT_COMMIT)' \
	-X '$(PKG_VERSION).GitVersion=$(GIT_VERSION)' \
	-X '$(PKG_VERSION).BuildDate=$(BUILD_DATE)'"

# For ultra-minimal size (~3.5MB), use:
# LDFLAGS_EXTRA=-ldflags "-s -w -extldflags=-static"

build:
	CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o $(BINARY_NAME) .

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
