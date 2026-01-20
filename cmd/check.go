package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Show current xrandr monitor layout",
	Long: `Display the current monitor configuration including active displays,
their resolutions, and which display is set as primary.`,
	Example: `  dmon check`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		layout, err := svc.CheckDisplays(getContext())
		if err != nil {
			return fmt.Errorf("failed to check displays: %w", err)
		}

		fmt.Println("Current Display Layout:")
		fmt.Println()

		activeCount := 0
		for _, d := range layout.Displays {
			if !d.Connected || d.CurrentMode == nil {
				continue
			}

			activeCount++
			primaryMarker := ""
			if d.ID == layout.Primary {
				primaryMarker = " [PRIMARY]"
			}

			fmt.Printf("▸ %s (%s)%s\n", d.ID, d.Type, primaryMarker)
			fmt.Printf("  Resolution: %s\n", d.CurrentMode)
			fmt.Println()
		}

		if activeCount == 0 {
			fmt.Println("No active displays found")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
