package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"password-manager-cli/internal/vault"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [output.json]",
	Short: "Export the vault to an unencrypted JSON file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputFile := args[0]
		withVault(func(v *vault.Vault, masterPw []byte, path string) {
			jsonData, err := json.MarshalIndent(v.Entries, "", "  ")
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
		})
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
