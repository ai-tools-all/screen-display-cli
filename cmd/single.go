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
		if err := svc.SetSingleDisplay(getContext()); err != nil {
			return fmt.Errorf("single display setup failed: %w", err)
		}

		fmt.Println("✓ Single display mode (internal only)")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(singleCmd)
}
