package xrandr

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/abhishek/dmon-cli/internal/models"
	"github.com/sirupsen/logrus"
)

type Backend struct {
	logger *logrus.Logger
}

func NewBackend(logger *logrus.Logger) *Backend {
	return &Backend{
		logger: logger,
	}
}

var (
	displayLineRegex = regexp.MustCompile(`^(\S+)\s+(connected|disconnected)`)
	modeLineRegex    = regexp.MustCompile(`^\s+(\d+)x(\d+)\s+([0-9.]+)([*+\s]*)`)
	internalPatterns = []string{"eDP", "LVDS"}
)

func (b *Backend) DetectDisplays(ctx context.Context) ([]models.Display, error) {
	b.logger.Debug("Detecting displays via xrandr")

	cmd := exec.CommandContext(ctx, "xrandr", "--query")
	output, err := cmd.Output()
	if err != nil {
		b.logger.WithError(err).Error("Failed to execute xrandr")
		return nil, fmt.Errorf("xrandr command failed: %w", err)
	}

	displays, err := b.parseXrandrOutput(string(output))
	if err != nil {
		b.logger.WithError(err).Error("Failed to parse xrandr output")
		return nil, err
	}

	connected := 0
	for _, d := range displays {
		if d.Connected {
			connected++
		}
	}

	b.logger.WithFields(logrus.Fields{
		"total":     len(displays),
		"connected": connected,
	}).Info("Displays detected")

	for _, d := range displays {
		b.logger.WithFields(logrus.Fields{
			"id":        d.ID,
			"type":      d.Type,
			"connected": d.Connected,
			"modes":     len(d.Modes),
		}).Debug("Display details")
	}

	return displays, nil
}

func (b *Backend) parseXrandrOutput(output string) ([]models.Display, error) {
	var displays []models.Display
	var currentDisplay *models.Display

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		if matches := displayLineRegex.FindStringSubmatch(line); matches != nil {
			if currentDisplay != nil {
				displays = append(displays, *currentDisplay)
			}

			displayID := matches[1]
			connected := matches[2] == "connected"

			currentDisplay = &models.Display{
				ID:        displayID,
				Type:      b.identifyDisplayType(displayID),
				Connected: connected,
				Modes:     []models.Mode{},
			}
			continue
		}

		if currentDisplay != nil && currentDisplay.Connected {
			if matches := modeLineRegex.FindStringSubmatch(line); matches != nil {
				width, _ := strconv.Atoi(matches[1])
				height, _ := strconv.Atoi(matches[2])
				rate, _ := strconv.ParseFloat(matches[3], 64)
				flags := matches[4]

				mode := models.Mode{
					Width:     width,
					Height:    height,
					Rate:      rate,
					Current:   strings.Contains(flags, "*"),
					Preferred: strings.Contains(flags, "+"),
				}

				currentDisplay.Modes = append(currentDisplay.Modes, mode)

				if mode.Current {
					modeCopy := mode
					currentDisplay.CurrentMode = &modeCopy
				}
			}
		}
	}

	if currentDisplay != nil {
		displays = append(displays, *currentDisplay)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return displays, nil
}

func (b *Backend) identifyDisplayType(displayID string) models.DisplayType {
	for _, pattern := range internalPatterns {
		if strings.HasPrefix(displayID, pattern) {
			return models.Internal
		}
	}
	return models.External
}

func (b *Backend) Configure(ctx context.Context, config models.DisplayConfig, displays []models.Display) (*models.ConfigResult, error) {
	b.logger.WithFields(logrus.Fields{
		"target":   config.Target,
		"mode":     config.Mode,
		"position": config.Position,
	}).Info("Configuring displays")

	internal, externals := b.categorizeDisplays(displays)

	if internal == nil {
		return nil, fmt.Errorf("no internal display found")
	}

	var args []string
	var configuredDisplays []models.ConfiguredDisplay

	switch config.Target {
	case models.TargetInternal:
		res := b.getResolution(internal, config.Mode)
		args = b.buildInternalOnlyConfig(internal, config.Mode)
		configuredDisplays = []models.ConfiguredDisplay{
			{
				ID:         internal.ID,
				Type:       internal.Type,
				Resolution: res,
				Active:     true,
			},
		}

	case models.TargetExternal:
		if len(externals) == 0 {
			return nil, fmt.Errorf("no external displays found. Try 'dmon list' to see available displays")
		}
		res := b.getResolution(externals[0], config.Mode)
		args = b.buildExternalOnlyConfig(internal, externals[0], config.Mode)
		configuredDisplays = []models.ConfiguredDisplay{
			{
				ID:         externals[0].ID,
				Type:       externals[0].Type,
				Resolution: res,
				Active:     true,
			},
			{
				ID:         internal.ID,
				Type:       internal.Type,
				Resolution: "",
				Active:     false,
			},
		}

	case models.TargetBoth:
		if len(externals) == 0 {
			return nil, fmt.Errorf("no external displays found. Try 'dmon list' to see available displays")
		}
		internalRes := b.getResolution(internal, config.Mode)
		externalRes := b.getResolution(externals[0], config.Mode)
		args = b.buildDualConfig(internal, externals[0], config.Mode, config.Position)
		configuredDisplays = []models.ConfiguredDisplay{
			{
				ID:         externals[0].ID,
				Type:       externals[0].Type,
				Resolution: externalRes,
				Active:     true,
			},
			{
				ID:         internal.ID,
				Type:       internal.Type,
				Resolution: internalRes,
				Active:     true,
			},
		}
	}

	b.logger.WithField("args", strings.Join(args, " ")).Debug("Executing xrandr command")

	cmd := exec.CommandContext(ctx, "xrandr", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		b.logger.WithFields(logrus.Fields{
			"error":  err,
			"output": string(output),
		}).Error("xrandr configuration failed")
		return nil, fmt.Errorf("xrandr failed: %w\nOutput: %s", err, string(output))
	}

	b.logger.Info("Display configuration applied successfully")

	result := &models.ConfigResult{
		Displays: configuredDisplays,
		Config:   config,
	}

	return result, nil
}

func (b *Backend) categorizeDisplays(displays []models.Display) (*models.Display, []*models.Display) {
	var internal *models.Display
	var externals []*models.Display

	for i := range displays {
		if !displays[i].Connected {
			continue
		}
		if displays[i].Type == models.Internal {
			internal = &displays[i]
		} else {
			externals = append(externals, &displays[i])
		}
	}

	return internal, externals
}

func (b *Backend) buildInternalOnlyConfig(internal *models.Display, mode models.ResolutionMode) []string {
	res := b.getResolution(internal, mode)
	return []string{
		"--output", internal.ID,
		"--mode", res,
		"--primary",
	}
}

func (b *Backend) buildExternalOnlyConfig(internal *models.Display, external *models.Display, mode models.ResolutionMode) []string {
	res := b.getResolution(external, mode)
	return []string{
		"--output", external.ID,
		"--mode", res,
		"--primary",
		"--output", internal.ID,
		"--off",
	}
}

func (b *Backend) buildDualConfig(internal *models.Display, external *models.Display, mode models.ResolutionMode, pos models.Position) []string {
	internalRes := b.getResolution(internal, mode)
	externalRes := b.getResolution(external, mode)

	args := []string{
		"--output", external.ID,
		"--mode", externalRes,
		"--primary",
		"--output", internal.ID,
		"--mode", internalRes,
	}

	switch pos {
	case models.PositionLeft:
		args = append(args, "--left-of", external.ID)
	case models.PositionRight:
		args = append(args, "--right-of", external.ID)
	case models.PositionAbove:
		args = append(args, "--above", external.ID)
	case models.PositionBelow:
		args = append(args, "--below", external.ID)
	default:
		args = append(args, "--right-of", external.ID)
	}

	return args
}

func (b *Backend) getResolution(display *models.Display, mode models.ResolutionMode) string {
	switch mode {
	case models.ModeNormal:
		if display.Type == models.Internal {
			return b.findClosestMode(display, 1920, 1200)
		}
		return b.findClosestMode(display, 1920, 1080)

	case models.ModeZoom:
		if display.Type == models.Internal {
			return b.findClosestMode(display, 1600, 1000)
		}
		return b.findClosestMode(display, 1280, 720)

	case models.ModeNative:
		return b.findNativeMode(display)

	default:
		return b.findNativeMode(display)
	}
}

func (b *Backend) findClosestMode(display *models.Display, targetWidth, targetHeight int) string {
	for _, mode := range display.Modes {
		if mode.Width == targetWidth && mode.Height == targetHeight {
			return fmt.Sprintf("%dx%d", mode.Width, mode.Height)
		}
	}

	if display.CurrentMode != nil {
		b.logger.WithFields(logrus.Fields{
			"display": display.ID,
			"target":  fmt.Sprintf("%dx%d", targetWidth, targetHeight),
		}).Warn("Target resolution not available, using current mode")
		return fmt.Sprintf("%dx%d", display.CurrentMode.Width, display.CurrentMode.Height)
	}

	return b.findNativeMode(display)
}

func (b *Backend) findNativeMode(display *models.Display) string {
	var best models.Mode
	maxPixels := 0

	for _, mode := range display.Modes {
		pixels := mode.Width * mode.Height
		if pixels > maxPixels {
			maxPixels = pixels
			best = mode
		}
	}

	if maxPixels > 0 {
		return fmt.Sprintf("%dx%d", best.Width, best.Height)
	}

	return "auto"
}

func (b *Backend) GetCurrentLayout(ctx context.Context) (*models.Layout, error) {
	b.logger.Debug("Getting current layout")

	displays, err := b.DetectDisplays(ctx)
	if err != nil {
		return nil, err
	}

	layout := &models.Layout{
		Displays: displays,
	}

	for _, d := range displays {
		if d.Connected && d.CurrentMode != nil {
			layout.Primary = d.ID
			break
		}
	}

	return layout, nil
}

func (b *Backend) GetSupportedModes(ctx context.Context, displayID string) ([]models.Mode, error) {
	b.logger.WithField("display", displayID).Debug("Getting supported modes")

	displays, err := b.DetectDisplays(ctx)
	if err != nil {
		return nil, err
	}

	for _, d := range displays {
		if d.ID == displayID {
			return d.Modes, nil
		}
	}

	return nil, fmt.Errorf("display %s not found", displayID)
}
