package cmd

import (
	"bytes"
	"fmt"
	"os"

	"password-manager-cli/internal/vault"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new password vault",
	Run: func(cmd *cobra.Command, args []string) {
		path := getVaultPath()

		// Check if vault already exists
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Vault already exists at", path)
			return
		}

		fmt.Println("Initializing new vault...")
		pw1, err := promptPassword("Enter Master Password: ")
		if err != nil || len(pw1) == 0 {
			fmt.Println("Error reading password")
			return
		}
		defer clear(pw1)

		// Validate Master Password strength
		if err := vault.ValidateMasterPassword(pw1); err != nil {
			fmt.Println("Weak Master Password:", err)
			return
		}

		pw2, err := promptPassword("Confirm Master Password: ")
		if err != nil {
			fmt.Println("Error reading password confirmation")
			return
		}
		defer clear(pw2)

		if !bytes.Equal(pw1, pw2) {
			fmt.Println("Passwords do not match!")
			return
		}

		// Generate new salt and create empty vault
		salt, err := vault.GenerateSalt(16)
		if err != nil {
			fmt.Println("Error generating salt:", err)
			return
		}

		v := &vault.Vault{
			Salt:    salt,
			Entries: make(map[string]vault.Entry),
		}

		err = vault.SaveVault(path, pw1, v)
		if err != nil {
			fmt.Println("Error saving vault:", err)
			return
		}

		fmt.Println("Vault successfully created at", path)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
