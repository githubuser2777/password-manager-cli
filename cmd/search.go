package cmd

import (
	"fmt"
	"password-manager-cli/internal/vault"
	"strings"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Search for a credential in the vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := strings.ToLower(args[0])
		withVault(func(v *vault.Vault, masterPw []byte, path string) {
			fmt.Printf("\nSearch results for '%s':\n", args[0])
			fmt.Println("----------------------------------------")

			found := false
			for service, entry := range v.Entries {
				if strings.Contains(strings.ToLower(service), keyword) || strings.Contains(strings.ToLower(entry.Username), keyword) {
					fmt.Printf("- %s (Username: %s)\n", service, entry.Username)
					found = true
				}
			}

			if !found {
				fmt.Println("No matches found.")
			}
			fmt.Println("----------------------------------------")
		})
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
