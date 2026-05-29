package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"password-manager-cli/internal/storage"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved services",
	Run: func(cmd *cobra.Command, args []string) {
		path := getVaultPath()

		masterPw, err := promptPassword("Master Password: ")
		if err != nil {
			return
		}

		vault, err := storage.LoadVault(path, masterPw)
		if err != nil {
			fmt.Println("Failed to open vault:", err)
			return
		}

		if len(vault.Entries) == 0 {
			fmt.Println("Vault is empty.")
			return
		}

		fmt.Println("Saved services:")
		for service, entry := range vault.Entries {
			fmt.Printf("- %s (Username: %s)\n", service, entry.Username)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
