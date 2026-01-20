package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var singleCmd = &cobra.Command{
	Use:   "single",
	Short: "Internal display only (disable external)",
	Long: `Switch to single display mode using only the internal display.
External displays will be disabled.`,
	Example: `  dmon single`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := svc.SetSingleDisplay(getContext())
		if err != nil {
			return fmt.Errorf("single display setup failed: %w", err)
		}

		fmt.Println("✓ Single display mode (internal only)\n")

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
	rootCmd.AddCommand(singleCmd)
}
