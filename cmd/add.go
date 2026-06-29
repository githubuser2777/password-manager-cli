package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"password-manager-cli/internal/vault"

	"github.com/spf13/cobra"
)

var generateFlag bool

var addCmd = &cobra.Command{
	Use:   "add [service]",
	Short: "Add a new credential to the vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]

		withVaultMutate(func(v *vault.Vault, masterPw []byte, path string) bool {
			if _, exists := v.Entries[service]; exists {
				fmt.Printf("Service '%s' already exists in the vault.\n", service)
				return false
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Username: ")
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)

			var password string
			if generateFlag {
				var err error
				password, err = vault.GenerateRandomPassword(16, true)
				if err != nil {
					fmt.Println("Error generating password:", err)
					return false
				}
				fmt.Println("Generated Password:", password)
			} else {
				rawPw, err := promptPassword("Password: ")
				if err != nil {
					return false
				}
				password = string(rawPw)
				clear(rawPw)
			}

			fmt.Print("Notes (optional): ")
			notes, _ := reader.ReadString('\n')
			notes = strings.TrimSpace(notes)

			now := time.Now().Format(time.RFC3339)
			v.Entries[service] = vault.Entry{
				Username:  username,
				Password:  password,
				Notes:     notes,
				CreatedAt: now,
				UpdatedAt: now,
			}

			fmt.Printf("Successfully added '%s' to the vault.\n", service)
			return true
		})
	},
}

func init() {
	addCmd.Flags().BoolVarP(&generateFlag, "generate", "g", false, "Auto-generate a secure random password")
	rootCmd.AddCommand(addCmd)
}
