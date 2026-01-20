package service

import (
	"context"
	"fmt"

	"github.com/abhishek/dmon-cli/internal/adapter"
	"github.com/abhishek/dmon-cli/internal/models"
	"github.com/sirupsen/logrus"
)

type DisplayService struct {
	backend adapter.DisplayBackend
	logger  *logrus.Logger
}

func New(backend adapter.DisplayBackend, logger *logrus.Logger) *DisplayService {
	return &DisplayService{
		backend: backend,
		logger:  logger,
	}
}

func (s *DisplayService) SetupDual(ctx context.Context, mode models.ResolutionMode) (*models.ConfigResult, error) {
	s.logger.WithField("mode", mode).Info("Setting up dual display")

	displays, err := s.backend.DetectDisplays(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to detect displays: %w", err)
	}

	config := models.DisplayConfig{
		Target:   models.TargetBoth,
		Mode:     mode,
		Position: models.PositionRight,
	}

	result, err := s.backend.Configure(ctx, config, displays)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *DisplayService) SetDisplay(ctx context.Context, target models.Target, mode models.ResolutionMode, position models.Position) (*models.ConfigResult, error) {
	s.logger.WithFields(logrus.Fields{
		"target":   target,
		"mode":     mode,
		"position": position,
	}).Info("Configuring display")

	displays, err := s.backend.DetectDisplays(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to detect displays: %w", err)
	}

	config := models.DisplayConfig{
		Target:   target,
		Mode:     mode,
		Position: position,
	}

	result, err := s.backend.Configure(ctx, config, displays)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *DisplayService) SetSingleDisplay(ctx context.Context) (*models.ConfigResult, error) {
	s.logger.Info("Setting up single display (internal only)")

	displays, err := s.backend.DetectDisplays(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to detect displays: %w", err)
	}

	config := models.DisplayConfig{
		Target:   models.TargetInternal,
		Mode:     models.ModePreset,
		Position: models.PositionNone,
	}

	result, err := s.backend.Configure(ctx, config, displays)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *DisplayService) ListDisplays(ctx context.Context) ([]models.Display, error) {
	s.logger.Info("Listing displays")

	displays, err := s.backend.DetectDisplays(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to detect displays: %w", err)
	}

	return displays, nil
}

func (s *DisplayService) DetectDisplays(ctx context.Context) ([]models.Display, error) {
	s.logger.Info("Re-detecting displays")

	displays, err := s.backend.DetectDisplays(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to detect displays: %w", err)
	}

	s.logger.WithField("count", len(displays)).Info("Display detection complete")
	return displays, nil
}

func (s *DisplayService) CheckDisplays(ctx context.Context) (*models.Layout, error) {
	s.logger.Info("Checking current display layout")

	layout, err := s.backend.GetCurrentLayout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current layout: %w", err)
	}

	return layout, nil
}
