# dmon - Display Monitor CLI

Golang CLI tool for managing display configurations on Linux via xrandr.

## Usage

```
dmon is a CLI tool for managing display configurations on Linux.
It provides a simple interface to xrandr for common display management tasks.

Usage:
  dmon [command]

Available Commands:
  check       Show current xrandr monitor layout
  completion  Generate the autocompletion script for the specified shell
  detect      Re-scan and update display inventory
  dual        Quick dual-display setup (external primary, internal right)
  help        Help about any command
  list        Show all connected displays with available modes
  set         Full control over display configuration
  single      Internal display only (disable external)

Flags:
  -h, --help      help for dmon
  -v, --verbose   Show detailed output and xrandr commands
      --version   version for dmon

Use "dmon [command] --help" for more information about a command.
```

## Features

- **Quick dual-display setup** - External primary, internal positioned
- **Flexible configuration** - Full control over target, mode, and position
- **Multiple resolution modes** - Preset, low, highest available
- **Verbose logging** - Human-readable stdout + structured JSON logs
- **Adapter pattern** - Ready for future backends (Wayland, etc.)
- **Display detection** - Re-scan for hot-plugged monitors
- **Shell completion** - Bash, Zsh, Fish, PowerShell support

## Installation

```bash
# Build from source
make build
# or: go build -o dmon .

# Install to local bin
make install
# or: cp dmon ~/.local/bin/

# Or system-wide
sudo cp dmon /usr/local/bin/
```

## Quick Start

```bash
# Dual display (external primary, internal right)
dmon dual

# Dual with reduced resolution
dmon dual low

# Dual with highest available resolution
dmon dual highest

# Internal display only
dmon single

# List available displays
dmon list

# Check current layout
dmon check

# Detect displays after plugging/unplugging
dmon detect
```

## Commands

### `dmon dual [mode]`
Configure dual-display mode with external monitor as primary and internal display positioned to the right.

**Modes:**
- `preset` (default) - 1920x1200 internal, 1920x1080 external
- `low` - 1600x1000 internal, 1280x720 external
- `highest` - Highest available resolution for each display

**Examples:**
```bash
dmon dual              # Uses preset mode
dmon dual low          # Uses low resolution
dmon dual highest      # Uses native resolution
```

### `dmon set <target> <mode> [position]`
Configure displays with complete control over target, mode, and positioning.

**Targets:**
- `internal` (i) - Internal display only
- `external` (e) - External display only
- `both` (b) - Both displays

**Modes:**
- `preset` (p) - 1920x1200 internal, 1920x1080 external
- `low` (l) - 1600x1000 internal, 1280x720 external
- `highest` (h) - Highest available resolution

**Positions** (optional, for 'both' target):
- `left` (l) - Internal display to the left of external
- `right` (r) - Internal display to the right of external (default)
- `above` (a) - Internal display above external
- `below` (b) - Internal display below external

**Examples:**
```bash
dmon set internal highest
dmon set external low
dmon set both preset left
dmon set i l           # Internal low (short form)
dmon set e h           # External highest (short form)
dmon set b p r         # Both preset right (short form)
```

### `dmon single`
Switch to single display mode using only the internal display. External displays will be disabled.

**Examples:**
```bash
dmon single
```

### `dmon list`
Display a list of all connected displays along with their supported resolutions. Shows which mode is currently active and which is the preferred mode.

**Examples:**
```bash
dmon list
```

### `dmon check`
Display the current monitor configuration including active displays, their resolutions, and which display is set as primary.

**Examples:**
```bash
dmon check
```

### `dmon detect`
Force a re-scan of connected displays and report what was found. Useful after plugging/unplugging external monitors.

**Examples:**
```bash
dmon detect
```

## Global Flags

- `-h, --help` - Show help information
- `-v, --verbose` - Show detailed output and xrandr commands executed
- `--version` - Display version information

## Resolution Modes Reference

| Mode | Internal | External |
|------|----------|----------|
| **preset** | 1920x1200 | 1920x1080 |
| **low** | 1600x1000 | 1280x720 |
| **highest** | Native max | Native max |

## Positioning Reference

| Position | Layout |
|----------|--------|
| **left** | Internal ← External |
| **right** | Internal → External |
| **above** | Internal ↑ External |
| **below** | Internal ↓ External |

## Usage Examples

### Basic scenarios
```bash
# Quick setup: external primary, internal to the right
dmon dual

# Use reduced resolution for performance
dmon dual low

# Switch to highest available resolutions
dmon dual highest

# Single screen mode
dmon single
```

### Advanced configuration
```bash
# Internal display in highest mode
dmon set internal highest

# External display in low mode
dmon set external low

# Both displays, preset mode, internal to the left
dmon set both preset left

# Short form equivalent
dmon set b p l
```

### Monitoring and troubleshooting
```bash
# View current configuration
dmon check

# List all connected displays and modes
dmon list

# Detect displays after hot-plugging
dmon detect

# See verbose output and xrandr commands
dmon -v dual
dmon --verbose set both preset
```

## Logs

- **Stdout**: Human-readable, colored output (INFO level by default, DEBUG with `-v`)
- **File**: `~/.local/share/dmon/dmon.log` (structured JSON format, always DEBUG level)

## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md) for design details, interfaces, and extending with new backends.

## Requirements

- Linux with X11
- xrandr installed
- Go 1.21+ (for building)

## License

MIT
