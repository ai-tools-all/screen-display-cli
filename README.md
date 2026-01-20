# dmon - Display Monitor CLI

Golang CLI tool for managing display configurations on Linux via xrandr.

## Features

- **Quick dual-display setup** - External primary, internal positioned
- **Flexible configuration** - Full control over target, mode, and position
- **Multiple resolution modes** - Normal, zoom, native
- **Verbose logging** - Human-readable stdout + structured JSON logs
- **Adapter pattern** - Ready for future backends (Wayland, etc.)

## Installation

```bash
# Build from source
go build -o dmon .

# Install to local bin
cp dmon ~/.local/bin/

# Or system-wide
sudo cp dmon /usr/local/bin/
```

## Quick Start

```bash
# Dual display (external primary, internal right)
dmon dual

# Dual with zoomed resolution
dmon dual zoom

# Internal display only
dmon single

# List available displays
dmon list

# Check current layout
dmon check
```

## Commands

| Command | Description |
|---------|-------------|
| `dmon dual [mode]` | Quick dual-display (external primary, internal right) |
| `dmon set <target> <mode> [position]` | Full control configuration |
| `dmon single` | Internal display only |
| `dmon list` | Show connected displays with modes |
| `dmon check` | Show current layout |
| `dmon detect` | Re-scan displays |

## Resolution Modes

- **normal** (n) - 1920x1200 internal, 1920x1080 external
- **zoom** (z) - 1600x1000 internal, 1280x720 external
- **native** (max) - Highest available resolution

## Positioning (for dual mode)

- **left** (l) - Internal left of external
- **right** (r) - Internal right of external
- **above** (a) - Internal above external
- **below** (b) - Internal below external

## Examples

```bash
# Advanced control
dmon set internal native left
dmon set external zoom above
dmon set both normal right

# Short form
dmon set i z l    # Internal zoom left
dmon set e n r    # External normal right

# Verbose mode
dmon -v dual zoom
dmon --verbose set both native
```

## Logs

- **Stdout**: Human-readable, colored (INFO level, DEBUG with `-v`)
- **File**: `~/.local/share/dmon/dmon.log` (structured JSON, always DEBUG)

## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md) for design details, interfaces, and extending with new backends.

## Requirements

- Linux with X11
- xrandr installed
- Go 1.21+ (for building)

## License

MIT
