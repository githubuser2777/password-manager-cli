package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"password-manager-cli/internal/core"
	"password-manager-cli/internal/crypto"
	"password-manager-cli/internal/storage"
)

var generateFlag bool

var addCmd = &cobra.Command{
	Use:   "add [service]",
	Short: "Add a new credential to the vault",
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

		if _, exists := vault.Entries[service]; exists {
			fmt.Printf("Service '%s' already exists in the vault.\n", service)
			return
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		var password string
		if generateFlag {
			password, err = crypto.GenerateRandomPassword(16, true)
			if err != nil {
				fmt.Println("Error generating password:", err)
				return
			}
			fmt.Println("Generated Password:", password)
		} else {
			password, _ = promptPassword("Password: ")
		}

		fmt.Print("Notes (optional): ")
		notes, _ := reader.ReadString('\n')
		notes = strings.TrimSpace(notes)

		now := time.Now().Format(time.RFC3339)
		vault.Entries[service] = core.Entry{
			Username:  username,
			Password:  password,
			Notes:     notes,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := storage.SaveVault(path, masterPw, vault); err != nil {
			fmt.Println("Failed to save vault:", err)
			return
		}

		fmt.Printf("Successfully added '%s' to the vault.\n", service)
	},
}

func init() {
	addCmd.Flags().BoolVarP(&generateFlag, "generate", "g", false, "Auto-generate a secure random password")
	rootCmd.AddCommand(addCmd)
}
