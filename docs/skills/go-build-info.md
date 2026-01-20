# Skill: Go Build Info & Versioning

## Description
This skill enables an agent to add semantic versioning, commit hashes, and build timestamps to a Golang binary using `ldflags`. It covers the creation of a version package, Makefile configuration, and Cobra CLI integration.

## Prerequisites
- A Go project (initialized with `go mod`).
- `make` (optional but recommended).
- `cobra` (optional, but standard for CLIs).
- Git repository.

## Implementation Steps

### 1. Create the Version Package
Create a lightweight, dependency-free package to hold the variables.
**Path:** `internal/version/version.go`

```go
package version

import (
	"fmt"
	"runtime"
)

// These variables are populated via -ldflags at build time.
// DO NOT modify these defaults manually.
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

### 2. Configure the Makefile
This is the **critical** step. You must construct the `LDFLAGS` to target the specific variables in the version package.

**Critical Configuration:**
Update `PKG_VERSION` to match the **full import path** of your version package.
Example: `github.com/org/repo/internal/version`

**Snippet for Makefile:**
```makefile
# --- Build Info Configuration ---
# TODO: Update this path to match your project's go.mod module path + /internal/version
PKG_VERSION=github.com/YOUR_ORG/YOUR_PROJECT/internal/version

# Dynamic variables
GIT_COMMIT=$(shell git rev-parse --short HEAD 2> /dev/null || echo "none")
GIT_VERSION=$(shell git describe --tags --always --dirty 2> /dev/null || echo "dev")
BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Linker flags
LDFLAGS=-ldflags "-s -w \
	-X '$(PKG_VERSION).GitCommit=$(GIT_COMMIT)' \
	-X '$(PKG_VERSION).GitVersion=$(GIT_VERSION)' \
	-X '$(PKG_VERSION).BuildDate=$(BUILD_DATE)'"

# --- Targets ---
.PHONY: build
build:
	go build $(LDFLAGS) -o bin/app ./cmd/main.go
```

### 3. Integrate with Cobra (CLI)
Wire the version info into your root command.

**Path:** `cmd/root.go`

```go
import (
    "fmt"
    "github.com/YOUR_ORG/YOUR_PROJECT/internal/version"
    // ...
)

func init() {
    // 1. Set the raw version string (used by --version flag)
    rootCmd.Version = version.GitVersion
    
    // 2. Customize the output template to show full info
    rootCmd.SetVersionTemplate(fmt.Sprintf("%s\n", version.Info()))

    // ... other flags
}
```

## Verification
Run the following to ensure variables are injected correctly:
```bash
make build
./bin/app --version
```

## Troubleshooting
- **Variables are "unknown" or "dev":**
    1. Check `PKG_VERSION` in the Makefile. It must match the output of `go list ./internal/version`.
    2. Ensure you are building via `make` (or passing the flags manually), not just `go build`.
