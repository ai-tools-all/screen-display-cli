# Display Management System - Requirements & Behavior

## High-Level Requirements

### 1. Dynamic Display Detection
**Requirement:** The system must automatically detect all connected displays on initialization without hardcoded values (except for internal display identification patterns).

**Behavior:**
- Scans system for all connected displays via `xrandr`
- Identifies internal displays by pattern: `eDP*`, `LVDS*`, `LVDS1`, or `eDP-1`
- Identifies all other connected displays as external displays
- Stores detection results in global variables `INTERNAL_DISPLAY` and `EXTERNAL_DISPLAYS` array

### 2. Flexible Resolution Management
**Requirement:** Support multiple resolution modes for both internal and external displays.

**Behavior:**
- **Normal Mode:** Default resolution (1920x1200 for internal, 1920x1080 for external)
- **Zoom Mode:** Reduced resolution for better visibility (1600x1000 internal, 1280x720 external)
- **Native/Max Mode:** Uses highest available resolution from display's supported modes

### 3. Intelligent External Display Handling
**Requirement:** Gracefully handle scenarios with and without external monitors.

**Behavior:**
- **With External Monitors:** Configures dual-display setup
- **Without External Monitors:** Falls back to internal display only (no errors)
- External display automatically set as primary when present
- Internal display positioned relative to external display

### 4. Positioning Control
**Requirement:** Support flexible display positioning for multi-monitor setups.

**Behavior:**
- **Left:** Internal display to the left of external
- **Right:** Internal display to the right of external (default for `setdualr`)
- **Above:** Internal display above external
- **Below:** Internal display below external
- **None:** No positional relationship (single display mode)

### 5. User-Friendly Interface
**Requirement:** Provide unified binary with intuitive subcommands and flags.

**Behavior:**

**Single Binary:** `dmon` (Display Monitor)

**Core Subcommands:**
- `dmon set <target> <mode> [position]` - Full control configuration
  - `<target>`: `internal` | `external` | `both` | `i` | `e` | `b`
  - `<mode>`: `normal` | `zoom` | `native` | `max` | `n` | `z`
  - `[position]`: `left` | `right` | `above` | `below` | `l` | `r` | `a` | `b` (optional)

- `dmon dual [mode]` - Quick dual-display (external primary, internal right)
  - `[mode]`: `normal` | `zoom` | `native` (default: normal)

- `dmon single` - Internal display only (disable external)

**Info & Management:**
- `dmon list` - Show all connected displays with available modes
- `dmon check` - Show current xrandr monitor layout
- `dmon detect` - Re-scan and update display inventory

**Global Flags:**
- `-v, --verbose` - Show executed xrandr commands
- `-h, --help` - Show help for any subcommand
- `--version` - Show version info

### 6. Error Handling & Feedback
**Requirement:** Provide clear feedback and fail gracefully.

**Behavior:**
- Detects when no displays are found and reports appropriately
- Detects when external displays are requested but not available
- Shows executed xrandr commands when `VERBOSE` is set
- Returns meaningful exit codes for scripting

### 7. Backward Compatibility
**Requirement:** Support migration from legacy shell functions.

**Behavior:**
- All legacy shell function capabilities mapped to subcommands
- Optional shell aliases can wrap `dmon` for legacy syntax:
  - `alias setdualr='dmon dual'`
  - `alias disables='dmon single'`
  - `alias listdisplays='dmon list'`
- Clean migration path from function-based to binary-based workflow

## Usage Examples

### Basic Setup
```bash
# Dual display with default settings
dmon dual

# Dual display with zoomed resolution
dmon dual zoom

# Dual display with native resolution
dmon dual native

# Single internal display only
dmon single
```

### Advanced Control
```bash
# Internal display at native resolution on the left
dmon set internal native left

# External display zoomed, positioned above internal
dmon set external zoom above

# Both displays at native resolution
dmon set both native

# Short form aliases
dmon set i z l    # Internal zoom on left
dmon set e n r    # External normal on right

# Check what displays are available
dmon list

# Verify current configuration
dmon check

# Re-detect displays
dmon detect
```

### Verbose Output
```bash
# Show xrandr commands being executed
dmon -v dual zoom
dmon --verbose set internal native left
```

### Help & Info
```bash
# General help
dmon --help

# Subcommand-specific help
dmon set --help
dmon dual --help

# Version info
dmon --version
```

## Technical Notes

### Binary Architecture
- Single compiled binary: `dmon` (Display Monitor)
- Stateless design: each invocation detects displays fresh
- Uses `xrandr` for all display operations (Linux/X11 specific)
- Exit codes: 0 (success), 1 (error), 2 (invalid args)

### Implementation Recommendations
- Language: Bash/Shell script (portable) or Go/Rust (performant)
- Configuration: Optional `~/.config/dmon/config` for defaults
- Caching: Optional `~/.cache/dmon/displays` for faster repeated calls
- Logging: Optional `--debug` flag for troubleshooting

### Installation
- Binary location: `/usr/local/bin/dmon` or `~/.local/bin/dmon`
- Optional shell completions: `/etc/bash_completion.d/dmon`
- Man page: `man dmon` for detailed documentation

## Command Reference

### Legacy Function → Binary Subcommand Mapping

| Legacy Function | New Binary Command | Notes |
|----------------|-------------------|-------|
| `setdualr` | `dmon dual` | Quick dual setup (default: normal) |
| `setdualr zoom` | `dmon dual zoom` | Dual with zoomed resolution |
| `setdualr native` | `dmon dual native` | Dual with native resolution |
| `setdisplay internal zoom left` | `dmon set internal zoom left` | Full control syntax |
| `setdisplay external normal right` | `dmon set external normal right` | Full control syntax |
| `disables` | `dmon single` | Internal display only |
| `listdisplays` | `dmon list` | Show available displays |
| `checkdisplays` | `dmon check` | Show current config |
| `detect_displays` | `dmon detect` | Re-scan displays |
| `dual` | `dmon dual` | Alias for setdualr |
| `dualz` | `dmon dual zoom` | Dual zoomed |
| `dualn` | `dmon dual native` | Dual native |
| `laptop` / `single` | `dmon single` | Single display mode |
| `sd i z` | `dmon set i z` | Shorthand still works |
| `sde n r` | `dmon set e n r` | Shorthand still works |

### Subcommand Summary

```
dmon dual [mode]                     → Quick dual-display setup
dmon set <target> <mode> [position]  → Full configuration control
dmon single                          → Internal display only
dmon list                            → Show available displays
dmon check                           → Show current layout
dmon detect                          → Re-scan displays
dmon --help                          → Show help
dmon --version                       → Show version
```
