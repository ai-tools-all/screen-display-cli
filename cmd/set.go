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
  preset, p    - Default resolution (1920x1200 internal, 1920x1080 external)
  low, l       - Reduced resolution (1600x1000 internal, 1280x720 external)
  highest, h   - Highest available resolution

Positions (optional, for 'both' target):
  left, l      - Internal display to the left of external
  right, r     - Internal display to the right of external (default)
  above, a     - Internal display above external
  below, b     - Internal display below external`,
	Example: `  dmon set internal highest
  dmon set external low
  dmon set both preset left
  dmon set i l
  dmon set e h`,
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

		result, err := svc.SetDisplay(getContext(), target, mode, position)
		if err != nil {
			return fmt.Errorf("display configuration failed: %w", err)
		}

		fmt.Printf("✓ Display configured (%s, %s", target, mode)
		if target == models.TargetBoth {
			fmt.Printf(", %s", position)
		}
		fmt.Println(")\n")

		fmt.Println("Configured displays:")
		for _, d := range result.Displays {
			if d.Active {
				fmt.Printf("  ▸ %s (%s) → %s\n", d.ID, d.Type, d.Resolution)
			} else {
				fmt.Printf("  ▸ %s (%s) → disabled\n", d.ID, d.Type)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
