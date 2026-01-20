package models

import "fmt"

type DisplayType int

const (
	Internal DisplayType = iota
	External
)

func (dt DisplayType) String() string {
	switch dt {
	case Internal:
		return "Internal"
	case External:
		return "External"
	default:
		return "Unknown"
	}
}

type Mode struct {
	Width     int
	Height    int
	Rate      float64
	Current   bool
	Preferred bool
}

func (m Mode) String() string {
	markers := ""
	if m.Current {
		markers += "*"
	}
	if m.Preferred {
		markers += "+"
	}
	return fmt.Sprintf("%dx%d@%.2fHz%s", m.Width, m.Height, m.Rate, markers)
}

type Display struct {
	ID          string
	Type        DisplayType
	Connected   bool
	Modes       []Mode
	CurrentMode *Mode
}

func (d Display) String() string {
	status := "disconnected"
	if d.Connected {
		status = "connected"
	}
	return fmt.Sprintf("%s (%s, %s)", d.ID, d.Type, status)
}

type Target int

const (
	TargetInternal Target = iota
	TargetExternal
	TargetBoth
)

func ParseTarget(s string) (Target, error) {
	switch s {
	case "internal", "i":
		return TargetInternal, nil
	case "external", "e":
		return TargetExternal, nil
	case "both", "b":
		return TargetBoth, nil
	default:
		return 0, fmt.Errorf("invalid target: %s (valid: internal/i, external/e, both/b)", s)
	}
}

func (t Target) String() string {
	switch t {
	case TargetInternal:
		return "internal"
	case TargetExternal:
		return "external"
	case TargetBoth:
		return "both"
	default:
		return "unknown"
	}
}

type ResolutionMode int

const (
	ModeNormal ResolutionMode = iota
	ModeZoom
	ModeNative
)

func ParseResolutionMode(s string) (ResolutionMode, error) {
	switch s {
	case "normal", "n":
		return ModeNormal, nil
	case "zoom", "z":
		return ModeZoom, nil
	case "native", "max":
		return ModeNative, nil
	default:
		return 0, fmt.Errorf("invalid mode: %s (valid: normal/n, zoom/z, native/max)", s)
	}
}

func (rm ResolutionMode) String() string {
	switch rm {
	case ModeNormal:
		return "normal"
	case ModeZoom:
		return "zoom"
	case ModeNative:
		return "native"
	default:
		return "unknown"
	}
}

type Position int

const (
	PositionNone Position = iota
	PositionLeft
	PositionRight
	PositionAbove
	PositionBelow
)

func ParsePosition(s string) (Position, error) {
	switch s {
	case "left", "l":
		return PositionLeft, nil
	case "right", "r":
		return PositionRight, nil
	case "above", "a":
		return PositionAbove, nil
	case "below", "b":
		return PositionBelow, nil
	case "none", "":
		return PositionNone, nil
	default:
		return 0, fmt.Errorf("invalid position: %s (valid: left/l, right/r, above/a, below/b)", s)
	}
}

func (p Position) String() string {
	switch p {
	case PositionLeft:
		return "left"
	case PositionRight:
		return "right"
	case PositionAbove:
		return "above"
	case PositionBelow:
		return "below"
	case PositionNone:
		return "none"
	default:
		return "unknown"
	}
}

type DisplayConfig struct {
	Target   Target
	Mode     ResolutionMode
	Position Position
}

type Layout struct {
	Displays []Display
	Primary  string
}
