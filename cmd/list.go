package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all connected displays with available modes",
	Long: `Display a list of all connected displays along with their supported resolutions.
Shows which mode is currently active and which is the preferred mode.`,
	Example: `  dmon list`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		displays, err := svc.ListDisplays(getContext())
		if err != nil {
			return fmt.Errorf("failed to list displays: %w", err)
		}

		if len(displays) == 0 {
			fmt.Println("No displays found")
			return nil
		}

		fmt.Printf("Found %d display(s):\n\n", len(displays))

		for _, d := range displays {
			if !d.Connected {
				continue
			}

			fmt.Printf("▸ %s (%s)\n", d.ID, d.Type)

			if len(d.Modes) > 0 {
				fmt.Println("  Available modes:")
				for _, m := range d.Modes {
					marker := "  "
					if m.Current {
						marker = "* "
					} else if m.Preferred {
						marker = "+ "
					}
					fmt.Printf("    %s %s\n", marker, m)
				}
			}
			fmt.Println()
		}

		fmt.Println("Legend: * = current, + = preferred")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
