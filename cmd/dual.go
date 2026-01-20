package cmd

import (
	"fmt"

	"github.com/abhishek/dmon-cli/internal/models"
	"github.com/spf13/cobra"
)

var dualCmd = &cobra.Command{
	Use:   "dual [mode]",
	Short: "Quick dual-display setup (external primary, internal right)",
	Long: `Configure dual-display mode with external monitor as primary and internal display positioned to the right.

Modes:
  preset  - Default resolution (1920x1200 internal, 1920x1080 external)
  low     - Reduced resolution (1600x1000 internal, 1280x720 external)
  highest - Highest available resolution for each display

If no mode is specified, 'preset' is used.`,
	Example: `  dmon dual
  dmon dual low
  dmon dual highest`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mode := models.ModePreset
		if len(args) > 0 {
			var err error
			mode, err = models.ParseResolutionMode(args[0])
			if err != nil {
				return err
			}
		}

		result, err := svc.SetupDual(getContext(), mode)
		if err != nil {
			return fmt.Errorf("dual display setup failed: %w", err)
		}

		fmt.Printf("✓ Dual display configured (%s mode)\n\n", mode)

		fmt.Println("Configured displays:")
		for _, d := range result.Displays {
			if d.Active {
				fmt.Printf("  ▸ %s (%s) → %s\n", d.ID, d.Type, d.Resolution)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dualCmd)
}
