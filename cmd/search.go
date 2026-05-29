package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"password-manager-cli/internal/crypto"
	"password-manager-cli/internal/storage"
)

var searchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Search for a credential in the vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := strings.ToLower(args[0])
		path := getVaultPath()

		masterPw, err := promptPassword("Master Password: ")
		if err != nil {
			return
		}
		defer crypto.ZeroBytes(masterPw)

		vault, err := storage.LoadVault(path, masterPw)
		if err != nil {
			fmt.Println("Failed to open vault:", err)
			return
		}

		fmt.Printf("\nSearch results for '%s':\n", args[0])
		fmt.Println("----------------------------------------")

		found := false
		for service, entry := range vault.Entries {
			if strings.Contains(strings.ToLower(service), keyword) || strings.Contains(strings.ToLower(entry.Username), keyword) {
				fmt.Printf("- %s (Username: %s)\n", service, entry.Username)
				found = true
			}
		}

		if !found {
			fmt.Println("No matches found.")
		}
		fmt.Println("----------------------------------------")
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
