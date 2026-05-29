package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"password-manager-cli/internal/core"
	"password-manager-cli/internal/storage"
)

var importCmd = &cobra.Command{
	Use:   "import [input.json]",
	Short: "Import credentials from an unencrypted JSON file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		path := getVaultPath()

		data, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Println("Failed to read import file:", err)
			return
		}

		var importedEntries map[string]core.Entry
		if err := json.Unmarshal(data, &importedEntries); err != nil {
			fmt.Println("Failed to parse JSON file (must be a map of services to entries):", err)
			return
		}

		masterPw, err := promptPassword("Master Password: ")
		if err != nil {
			return
		}

		vault, err := storage.LoadVault(path, masterPw)
		if err != nil {
			fmt.Println("Failed to open vault (make sure it's initialized first):", err)
			return
		}

		count := 0
		for k, v := range importedEntries {
			vault.Entries[k] = v
			count++
		}

		if err := storage.SaveVault(path, masterPw, vault); err != nil {
			fmt.Println("Failed to save vault:", err)
			return
		}

		fmt.Printf("Successfully imported %d credentials into the vault.\n", count)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
