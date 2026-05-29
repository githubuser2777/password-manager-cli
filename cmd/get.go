package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"password-manager-cli/internal/storage"
)

var copyFlag bool

var getCmd = &cobra.Command{
	Use:   "get [service]",
	Short: "Get a credential from the vault",
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
			fmt.Printf("Service '%s' not found.\n", service)
			return
		}

		fmt.Println("Username:", entry.Username)

		if copyFlag {
			if err := clipboard.WriteAll(entry.Password); err != nil {
				fmt.Println("Failed to copy to clipboard:", err)
			} else {
				fmt.Println("Password copied to clipboard!")
			}
		} else {
			fmt.Println("Password:", entry.Password)
		}
	},
}

func init() {
	getCmd.Flags().BoolVarP(&copyFlag, "copy", "c", false, "Copy password to clipboard instead of displaying it")
	rootCmd.AddCommand(getCmd)
}
