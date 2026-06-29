package cmd

import (
	"fmt"
	"os"

	"password-manager-cli/internal/tui"

	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the terminal user interface",
	Long:  `Launch the interactive Terminal User Interface (TUI) for passmgr.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := tui.StartApp(getVaultPath())
		if err != nil {
			fmt.Println("Error running TUI:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
