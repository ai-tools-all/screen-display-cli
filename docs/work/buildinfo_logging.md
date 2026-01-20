# Universal BuildInfo & Telemetry Pattern

This guide provides a drop-in pattern for adding versioning and dual-mode logging (Human-readable CLI + Structured JSONL) to any Go project.

## 1. The Version Package
Create a standalone package to hold build metadata. This isolates the logic and avoids circular dependencies.

**File:** `pkg/version/version.go`
```go
package version

import (
	"fmt"
	"runtime"
)

// These variables are populated via -ldflags at build time.
var (
	GitCommit  = "unknown"
	GitVersion = "v0.0.0-dev"
	BuildDate  = "unknown"
)

// Info returns a formatted string containing all build details.
func Info() string {
	return fmt.Sprintf("Version: %s\nCommit:  %s\nDate:    %s\nRuntime: %s/%s",
		GitVersion, GitCommit, BuildDate, runtime.GOOS, runtime.GOARCH)
}

// Map returns build details as a map for structured logging.
func Map() map[string]string {
	return map[string]string{
		"version": GitVersion,
		"commit":  GitCommit,
		"date":    BuildDate,
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
	}
}
```

## 2. The Build Command (Makefile)
Use this standard `Makefile` block to inject the variables automatically.

**File:** `Makefile`
```makefile
# Project config
BINARY_NAME=app
PKG_VERSION=github.com/yourusername/yourproject/pkg/version

# Dynamic variables
GIT_COMMIT=$(shell git rev-parse --short HEAD 2> /dev/null || echo "none")
GIT_VERSION=$(shell git describe --tags --always --dirty 2> /dev/null || echo "dev")
BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Linker flags
LDFLAGS=-ldflags "-s -w \
	-X '$(PKG_VERSION).GitCommit=$(GIT_COMMIT)' \
	-X '$(PKG_VERSION).GitVersion=$(GIT_VERSION)' \
	-X '$(PKG_VERSION).BuildDate=$(BUILD_DATE)'"

.PHONY: build
build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/main.go
```

## 3. Dual-Mode Logging (slog)
This implementation provides a "Fan-Out" logger that writes:
1.  **Stdout:** Human-friendly text (colored, simple).
2.  **File:** Machine-friendly JSONL (detailed, trace-ready).

**File:** `pkg/telemetry/logger.go`
```go
package telemetry

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// SetupLogger configures the global logger.
// logFile: path to the jsonl log file (e.g., "app.jsonl")
// verbose: if true, sets level to DEBUG, otherwise INFO.
func SetupLogger(logPath string, verbose bool) error {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	// 1. File Handler (JSONL)
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	jsonHandler := slog.NewJSONHandler(f, &slog.HandlerOptions{
		Level: level,
	})

	// 2. Stdout Handler (Text)
	// You can use a 3rd party library like 'lmittmann/tint' for colors, 
	// or the stdlib TextHandler.
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	// 3. Fan-out Handler
	multiHandler := &FanOutHandler{
		handlers: []slog.Handler{jsonHandler, textHandler},
	}

	logger := slog.New(multiHandler)
	slog.SetDefault(logger)
	
	return nil
}

// FanOutHandler broadcasts records to multiple handlers.
type FanOutHandler struct {
	handlers []slog.Handler
}

func (h *FanOutHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (h *FanOutHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			_ = handler.Handle(ctx, r)
		}
	}
	return nil
}

func (h *FanOutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &FanOutHandler{handlers: newHandlers}
}

func (h *FanOutHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &FanOutHandler{handlers: newHandlers}
}
```

## 4. Integration in Main
Wire it all together in your application entry point.

**File:** `cmd/main.go`
```go
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/yourusername/yourproject/pkg/telemetry"
	"github.com/yourusername/yourproject/pkg/version"
)

func main() {
	// Flags
	showVer := flag.Bool("version", false, "Show version info")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	logFile := flag.String("log", "app.jsonl", "Path to log file")
	flag.Parse()

	// 1. Handle Version
	if *showVer {
		fmt.Println(version.Info())
		os.Exit(0)
	}

	// 2. Setup Logging
	if err := telemetry.SetupLogger(*logFile, *verbose); err != nil {
		panic(err)
	}

	// 3. Log Startup
	slog.Info("Application starting", "build_info", version.Map())
	
	// App logic here...
}
```
