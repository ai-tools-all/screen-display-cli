package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Re-scan and update display inventory",
	Long: `Force a re-scan of connected displays and report what was found.
Useful after plugging/unplugging external monitors.`,
	Example: `  dmon detect`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		displays, err := svc.DetectDisplays(getContext())
		if err != nil {
			return fmt.Errorf("display detection failed: %w", err)
		}

		connectedCount := 0
		for _, d := range displays {
			if d.Connected {
				connectedCount++
			}
		}

		fmt.Printf("✓ Display scan complete\n")
		fmt.Printf("  Total displays: %d\n", len(displays))
		fmt.Printf("  Connected: %d\n", connectedCount)
		fmt.Println()

		if connectedCount > 0 {
			fmt.Println("Connected displays:")
			for _, d := range displays {
				if d.Connected {
					fmt.Printf("  ▸ %s (%s) - %d modes available\n", d.ID, d.Type, len(d.Modes))
				}
			}
			fmt.Println("\nUse 'dmon list' to see detailed mode information")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
}
