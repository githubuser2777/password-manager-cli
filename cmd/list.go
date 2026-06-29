package cmd

import (
	"fmt"
	"password-manager-cli/internal/vault"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved services",
	Run: func(cmd *cobra.Command, args []string) {
		withVault(func(v *vault.Vault, masterPw []byte, path string) {
			if len(v.Entries) == 0 {
				fmt.Println("Vault is empty.")
				return
			}

			fmt.Println("Saved services:")
			for service, entry := range v.Entries {
				fmt.Printf("- %s (Username: %s)\n", service, entry.Username)
			}
		})
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
