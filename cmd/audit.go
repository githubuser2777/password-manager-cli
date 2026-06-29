package cmd

import (
	"fmt"

	"password-manager-cli/internal/vault"

	"github.com/spf13/cobra"
)

var onlineFlag bool

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit passwords for strength and potential leaks",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		withVault(func(v *vault.Vault, masterPw []byte, path string) {
			fmt.Println("\nStarting Password Audit...")
			fmt.Println("----------------------------------------")

			report := vault.RunAudit(v, onlineFlag)

			for _, msg := range report.Messages {
				switch msg.Level {
				case "REUSED":
					fmt.Printf("[!] %s\n", msg.Message)
				case "WEAK":
					fmt.Printf("[!] %s\n", msg.Message)
				case "PWNED":
					fmt.Printf("[!!!] %s\n", msg.Message)
				case "ERROR":
					fmt.Printf("[-] %s\n", msg.Message)
				}
			}

			fmt.Println("----------------------------------------")
			fmt.Printf("Audit Complete! Total passwords checked: %d\n", report.TotalChecked)
			fmt.Printf("Weak: %d | Reused: %d | Pwned: %d\n", report.WeakCount, report.ReusedCount, report.PwnedCount)
			if !onlineFlag {
				fmt.Println("(Run with --online to check for data breaches via HaveIBeenPwned API)")
			}
		})
	},
}

func init() {
	auditCmd.Flags().BoolVarP(&onlineFlag, "online", "o", false, "Check HaveIBeenPwned API for data breaches (k-Anonymity)")
	rootCmd.AddCommand(auditCmd)
}
