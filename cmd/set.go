package cmd

import (
	"fmt"

	"github.com/abhishek/dmon-cli/internal/models"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set <target> <mode> [position]",
	Short: "Full control over display configuration",
	Long: `Configure displays with complete control over target, mode, and positioning.

Targets:
  internal, i  - Internal display only
  external, e  - External display only
  both, b      - Both displays

Modes:
  normal, n    - Default resolution (1920x1200 internal, 1920x1080 external)
  zoom, z      - Reduced resolution (1600x1000 internal, 1280x720 external)
  native, max  - Highest available resolution

Positions (optional, for 'both' target):
  left, l      - Internal display to the left of external
  right, r     - Internal display to the right of external (default)
  above, a     - Internal display above external
  below, b     - Internal display below external`,
	Example: `  dmon set internal native
  dmon set external zoom
  dmon set both normal left
  dmon set i z l
  dmon set e n r`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		target, err := models.ParseTarget(args[0])
		if err != nil {
			return err
		}

		mode, err := models.ParseResolutionMode(args[1])
		if err != nil {
			return err
		}

		position := models.PositionRight
		if len(args) > 2 {
			position, err = models.ParsePosition(args[2])
			if err != nil {
				return err
			}
		}

		if err := svc.SetDisplay(getContext(), target, mode, position); err != nil {
			return fmt.Errorf("display configuration failed: %w", err)
		}

		fmt.Printf("✓ Display configured (%s, %s", target, mode)
		if target == models.TargetBoth {
			fmt.Printf(", %s", position)
		}
		fmt.Println(")")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
