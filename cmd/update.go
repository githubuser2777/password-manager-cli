package cmd

import (
	"bufio"
	"fmt"
	"os"
	"password-manager-cli/internal/vault"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var updateGenerateFlag bool

var updateCmd = &cobra.Command{
	Use:   "update [service]",
	Short: "Update an existing credential in the vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
		withVaultMutate(func(v *vault.Vault, masterPw []byte, path string) bool {
			entry, exists := v.Entries[service]
			if !exists {
				fmt.Printf("Service '%s' not found in the vault.\n", service)
				return false
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
				var err error
				newPassword, err = vault.GenerateRandomPassword(16, true)
				if err != nil {
					fmt.Println("Error generating password:", err)
					return false
				}
				fmt.Println("Generated Password:", newPassword)
			} else {
				rawPw, err := promptPassword("Password [keep existing]: ")
				if err != nil {
					return false
				}
				if len(rawPw) == 0 {
					newPassword = entry.Password
				} else {
					newPassword = string(rawPw)
					clear(rawPw)
				}
			}

			fmt.Printf("Notes [%s]: ", entry.Notes)
			newNotes, _ := reader.ReadString('\n')
			newNotes = strings.TrimSpace(newNotes)
			if newNotes == "" {
				newNotes = entry.Notes
			}

			entry.Username = newUsername
			entry.Password = newPassword
			entry.Notes = newNotes
			entry.UpdatedAt = time.Now().Format(time.RFC3339)

			v.Entries[service] = entry
			fmt.Printf("Successfully updated '%s'.\n", service)
			return true
		})
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&updateGenerateFlag, "generate", "g", false, "Auto-generate a secure random password")
	rootCmd.AddCommand(updateCmd)
}
