package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"password-manager-cli/internal/crypto"
	"password-manager-cli/internal/storage"
)

var changepassCmd = &cobra.Command{
	Use:   "changepass",
	Short: "Change the Master Password",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		path := getVaultPath()

		oldPw, err := promptPassword("Current Master Password: ")
		if err != nil {
			return
		}
		defer crypto.ZeroBytes(oldPw)

		vault, err := storage.LoadVault(path, oldPw)
		if err != nil {
			fmt.Println("Failed to open vault (incorrect current password?):", err)
			return
		}

		newPw1, err := promptPassword("New Master Password: ")
		if err != nil {
			return
		}
		defer crypto.ZeroBytes(newPw1)

		// Validate Master Password strength
		if err := crypto.ValidateMasterPassword(newPw1); err != nil {
			fmt.Println("Weak Master Password:", err)
			return
		}

		newPw2, err := promptPassword("Confirm New Master Password: ")
		if err != nil {
			return
		}
		defer crypto.ZeroBytes(newPw2)

		if !bytes.Equal(newPw1, newPw2) {
			fmt.Println("Passwords do not match. Aborting.")
			return
		}

		if len(newPw1) == 0 {
			fmt.Println("Password cannot be empty.")
			return
		}

		if err := storage.SaveVault(path, newPw1, vault); err != nil {
			fmt.Println("Failed to save vault with new password:", err)
			return
		}

		fmt.Println("Master Password changed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(changepassCmd)
}
