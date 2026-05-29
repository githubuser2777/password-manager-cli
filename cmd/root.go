package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "passmgr",
	Short: "A secure password manager CLI",
	Long: `passmgr is a command-line password manager written in Go.
It uses AES-256-GCM for encryption and Argon2id for key derivation
to securely store your passwords in a local vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		tuiCmd.Run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
}
