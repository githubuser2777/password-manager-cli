package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"password-manager-cli/internal/storage"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [service]",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete a credential from the vault",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
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

		if _, exists := vault.Entries[service]; !exists {
			fmt.Printf("Service '%s' not found in the vault.\n", service)
			return
		}

		delete(vault.Entries, service)

		if err := storage.SaveVault(path, masterPw, vault); err != nil {
			fmt.Println("Failed to save vault:", err)
			return
		}

		fmt.Printf("Successfully deleted '%s' from the vault.\n", service)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
