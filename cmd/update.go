package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"password-manager-cli/internal/crypto"
	"password-manager-cli/internal/storage"
)

var updateGenerateFlag bool

var updateCmd = &cobra.Command{
	Use:   "update [service]",
	Short: "Update an existing credential in the vault",
	Args:  cobra.ExactArgs(1),
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

		entry, exists := vault.Entries[service]
		if !exists {
			fmt.Printf("Service '%s' not found in the vault.\n", service)
			return
		}

		fmt.Printf("Updating '%s'. Leave blank to keep existing value.\n", service)

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Username [%s]: ", entry.Username)
		newUsername, _ := reader.ReadString('\n')
		newUsername = strings.TrimSpace(newUsername)
		if newUsername == "" {
			newUsername = entry.Username
		}

		var newPassword string
		if updateGenerateFlag {
			newPassword, err = crypto.GenerateRandomPassword(16, true)
			if err != nil {
				fmt.Println("Error generating password:", err)
				return
			}
			fmt.Println("Generated Password:", newPassword)
		} else {
			newPassword, _ = promptPassword("Password [keep existing]: ")
			if newPassword == "" {
				newPassword = entry.Password
			}
		}

		fmt.Printf("Notes [%s]: ", entry.Notes)
		newNotes, _ := reader.ReadString('\n')
		newNotes = strings.TrimSpace(newNotes)
		// If user wants to clear notes, they might enter a space? We trim space. Let's just say empty means keep existing. 
		// If they want to clear, maybe they have to use a flag, but for now empty keeps existing.
		if newNotes == "" {
			newNotes = entry.Notes
		}

		entry.Username = newUsername
		entry.Password = newPassword
		entry.Notes = newNotes
		entry.UpdatedAt = time.Now().Format(time.RFC3339)

		vault.Entries[service] = entry

		if err := storage.SaveVault(path, masterPw, vault); err != nil {
			fmt.Println("Failed to save vault:", err)
			return
		}

		fmt.Printf("Successfully updated '%s'.\n", service)
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&updateGenerateFlag, "generate", "g", false, "Auto-generate a secure random password")
	rootCmd.AddCommand(updateCmd)
}
