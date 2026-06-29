package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"password-manager-cli/internal/vault"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [input.json]",
	Short: "Import credentials from an unencrypted JSON file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		
		data, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Println("Failed to read import file:", err)
			return
		}

		var importedEntries map[string]vault.Entry
		if err := json.Unmarshal(data, &importedEntries); err != nil {
			fmt.Println("Failed to parse JSON file (must be a map of services to entries):", err)
			return
		}

		withVaultMutate(func(v *vault.Vault, masterPw []byte, path string) bool {
			count := 0
			for k, vEntry := range importedEntries {
				v.Entries[k] = vEntry
				count++
			}
			fmt.Printf("Successfully imported %d credentials into the vault.\n", count)
			return count > 0
		})
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
