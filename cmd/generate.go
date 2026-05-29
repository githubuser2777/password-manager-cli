package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"password-manager-cli/internal/crypto"
)

var noSpecial bool

var generateCmd = &cobra.Command{
	Use:   "generate [length]",
	Short: "Generate a secure random password",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		length := 16
		if len(args) == 1 {
			if l, err := strconv.Atoi(args[0]); err == nil && l >= 4 {
				length = l
			} else {
				fmt.Println("Invalid length. Using default (16).")
			}
		}

		password, err := crypto.GenerateRandomPassword(length, !noSpecial)
		if err != nil {
			fmt.Println("Error generating password:", err)
			return
		}

		fmt.Println(password)
	},
}

func init() {
	generateCmd.Flags().BoolVarP(&noSpecial, "no-special", "n", false, "Exclude special characters")
	rootCmd.AddCommand(generateCmd)
}
