package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"password-manager-cli/internal/crypto"
	"password-manager-cli/internal/storage"
)

var exportCmd = &cobra.Command{
	Use:   "export [output.json]",
	Short: "Export the vault to an unencrypted JSON file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputFile := args[0]
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

		jsonData, err := json.MarshalIndent(vault.Entries, "", "  ")
		if err != nil {
			fmt.Println("Failed to marshal vault data:", err)
			return
		}

		if err := os.WriteFile(outputFile, jsonData, 0600); err != nil {
			fmt.Println("Failed to write export file:", err)
			return
		}

		fmt.Printf("Successfully exported vault to '%s'.\n", outputFile)
		fmt.Println("WARNING: This file is unencrypted. Please store it securely or delete it after use.")
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
