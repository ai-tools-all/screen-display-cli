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
