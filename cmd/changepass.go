package cmd

import (
	"bytes"
	"fmt"
	"password-manager-cli/internal/vault"

	"github.com/spf13/cobra"
)

var changepassCmd = &cobra.Command{
	Use:   "changepass",
	Short: "Change the Master Password",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		withVault(func(v *vault.Vault, oldPw []byte, path string) {
			newPw1, err := promptPassword("New Master Password: ")
			if err != nil {
				return
			}
			defer clear(newPw1)

			// Validate Master Password strength
			if err := vault.ValidateMasterPassword(newPw1); err != nil {
				fmt.Println("Weak Master Password:", err)
				return
			}

			newPw2, err := promptPassword("Confirm New Master Password: ")
			if err != nil {
				return
			}
			defer clear(newPw2)

			if !bytes.Equal(newPw1, newPw2) {
				fmt.Println("Passwords do not match. Aborting.")
				return
			}

			if len(newPw1) == 0 {
				fmt.Println("Password cannot be empty.")
				return
			}

			if err := vault.SaveVault(path, newPw1, v); err != nil {
				fmt.Println("Failed to save vault with new password:", err)
				return
			}

			fmt.Println("Master Password changed successfully!")
		})
	},
}

func init() {
	rootCmd.AddCommand(changepassCmd)
}
