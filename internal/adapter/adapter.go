package adapter

import (
	"context"

	"github.com/abhishek/dmon-cli/internal/models"
)

type DisplayDetector interface {
	DetectDisplays(ctx context.Context) ([]models.Display, error)
}

type DisplayConfigurator interface {
	Configure(ctx context.Context, config models.DisplayConfig, displays []models.Display) error
}

type DisplayQuerier interface {
	GetCurrentLayout(ctx context.Context) (*models.Layout, error)
	GetSupportedModes(ctx context.Context, displayID string) ([]models.Mode, error)
}

type DisplayBackend interface {
	DisplayDetector
	DisplayConfigurator
	DisplayQuerier
}
