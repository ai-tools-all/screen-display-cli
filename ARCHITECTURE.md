# dmon Architecture

## Overview

`dmon` (Display Monitor) is a Golang CLI tool for managing display configurations on Linux. It orchestrates xrandr with a clean adapter pattern, allowing future support for other display backends (Wayland, Windows, etc.).

## Design Principles

- **Composable Interfaces** - Small, focused interfaces that combine for flexibility
- **Adapter Pattern** - Abstract backend implementation from business logic
- **Stateless Execution** - No persistent state, fresh detection each invocation
- **Dual Logging** - Human-readable stdout + structured JSON file logs

## Architecture Layers

```
┌─────────────────────────────────────┐
│         CLI Layer (Cobra)           │
│  Commands: dual, set, single, etc.  │
└─────────────┬───────────────────────┘
              │
┌─────────────▼───────────────────────┐
│      Service Layer (Business)       │
│  Resolution mapping, orchestration  │
└─────────────┬───────────────────────┘
              │
┌─────────────▼───────────────────────┐
│    Adapter Layer (Interfaces)       │
│  DisplayDetector, Configurator,     │
│  Querier → DisplayBackend           │
└─────────────┬───────────────────────┘
              │
┌─────────────▼───────────────────────┐
│   Backend Implementation (xrandr)   │
│  Parses output, builds commands     │
└─────────────────────────────────────┘
```

## Project Structure

```
dmon-cli/
├── cmd/                    # Cobra commands (CLI interface)
│   ├── root.go            # Root command + global setup
│   ├── dual.go            # Quick dual-display setup
│   ├── set.go             # Full control configuration
│   ├── single.go          # Internal display only
│   ├── list.go            # Show available displays
│   ├── check.go           # Current layout status
│   └── detect.go          # Re-scan displays
│
├── internal/
│   ├── adapter/           # Backend interfaces
│   │   └── adapter.go     # DisplayBackend interface
│   │
│   ├── models/            # Data structures
│   │   └── types.go       # Display, Mode, Config types
│   │
│   ├── xrandr/            # xrandr backend implementation
│   │   └── xrandr.go      # Parse output, build commands
│   │
│   ├── service/           # Business logic
│   │   └── service.go     # Resolution mapping, orchestration
│   │
│   └── logger/            # Logging setup
│       └── logger.go      # Dual output (stdout + file)
│
├── main.go                # Entry point
├── go.mod                 # Go module definition
└── go.sum                 # Dependency checksums
```

## Key Interfaces

### DisplayBackend (Composable)

```go
type DisplayDetector interface {
    DetectDisplays(ctx context.Context) ([]Display, error)
}

type DisplayConfigurator interface {
    Configure(ctx context.Context, config DisplayConfig, displays []Display) error
}

type DisplayQuerier interface {
    GetCurrentLayout(ctx context.Context) (*Layout, error)
    GetSupportedModes(ctx context.Context, displayID string) ([]Mode, error)
}

type DisplayBackend interface {
    DisplayDetector
    DisplayConfigurator
    DisplayQuerier
}
```

## Data Flow

1. **CLI** → Cobra parses command + flags
2. **Service** → Detects displays via backend
3. **Service** → Maps mode (normal/zoom/native) to resolution
4. **Backend** → Builds xrandr command with proper args
5. **Backend** → Executes xrandr and captures output
6. **CLI** → Reports success/failure to user

## Logging Strategy

- **Stdout**: Human-readable, colored, INFO level (DEBUG with `-v`)
- **File**: `~/.local/share/dmon/dmon.log`, structured JSON, always DEBUG level
- **Library**: logrus with custom formatters + dual-output hook

## Adding New Backends

To add support for a new display system (e.g., Wayland):

1. Create `internal/wlrandr/wlrandr.go`
2. Implement `adapter.DisplayBackend` interface
3. Add backend selection logic in `cmd/root.go`
4. No changes needed to service layer or CLI commands

## Resolution Modes

| Mode   | Internal     | External     | Logic                    |
|--------|--------------|--------------|--------------------------|
| Normal | 1920x1200    | 1920x1080    | Find closest match       |
| Zoom   | 1600x1000    | 1280x720     | Find closest match       |
| Native | Highest      | Highest      | Max width × height       |

## Display Detection

- Internal displays: Match `eDP*` or `LVDS*` prefix patterns
- External displays: All others connected via HDMI/DP/VGA
- Detection: Parse `xrandr --query` output with regex
- Modes: Extract resolution, refresh rate, current/preferred flags

## Error Handling

- Balanced approach: User-friendly message + helpful hint
- Example: "External display not found. Try 'dmon list' to see available displays"
- Verbose mode: Shows full xrandr command + output on failure
- Exit codes: 0 (success), 1 (error), 2 (invalid args)

## Future Enhancements

- [ ] Config file support (`~/.config/dmon/config.toml`)
- [ ] Display profile saving/loading
- [ ] Wayland backend (wlr-randr)
- [ ] Brightness control
- [ ] Auto-switching on display connect/disconnect
- [ ] Shell completion scripts
- [ ] Man page generation
