package cmd

import (
	"fmt"
	"password-manager-cli/internal/vault"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [service]",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete a credential from the vault",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
		withVaultMutate(func(v *vault.Vault, masterPw []byte, path string) bool {
			if _, exists := v.Entries[service]; !exists {
				fmt.Printf("Service '%s' not found in the vault.\n", service)
				return false
			}

			delete(v.Entries, service)
			fmt.Printf("Successfully deleted '%s' from the vault.\n", service)
			return true
		})
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
