package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"password-manager-cli/internal/core"
	"password-manager-cli/internal/storage"
)

var onlineFlag bool

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit passwords for strength and potential leaks",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
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

		fmt.Println("\nStarting Password Audit...")
		fmt.Println("----------------------------------------")

		weakCount := 0
		reusedCount := 0
		pwnedCount := 0

		// Check for reuse
		pwMap := make(map[string][]string)
		for service, entry := range vault.Entries {
			pwMap[entry.Password] = append(pwMap[entry.Password], service)
		}

		for _, services := range pwMap {
			if len(services) > 1 {
				reusedCount++
				fmt.Printf("[!] REUSED: The password for %s is used across %d services.\n", services[0], len(services))
			}
		}

		for service, entry := range vault.Entries {
			// Basic strength check
			if len(entry.Password) < 8 {
				weakCount++
				fmt.Printf("[!] WEAK: Password for '%s' is too short (under 8 chars).\n", service)
			} else if !strings.ContainsAny(entry.Password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") ||
				!strings.ContainsAny(entry.Password, "0123456789") {
				weakCount++
				fmt.Printf("[!] WEAK: Password for '%s' lacks numbers or uppercase letters.\n", service)
			}

			// Online HIBP check
			if onlineFlag {
				pwned, err := core.CheckPwned(entry.Password)
				if err != nil {
					fmt.Printf("[-] HIBP Check failed for '%s': %v\n", service, err)
				} else if pwned {
					pwnedCount++
					fmt.Printf("[!!!] PWNED: Password for '%s' has been found in data breaches!\n", service)
				}
			}
		}

		fmt.Println("----------------------------------------")
		fmt.Printf("Audit Complete! Total passwords checked: %d\n", len(vault.Entries))
		fmt.Printf("Weak: %d | Reused: %d | Pwned: %d\n", weakCount, reusedCount, pwnedCount)
		if !onlineFlag {
			fmt.Println("(Run with --online to check for data breaches via HaveIBeenPwned API)")
		}
	},
}

func init() {
	auditCmd.Flags().BoolVarP(&onlineFlag, "online", "o", false, "Check HaveIBeenPwned API for data breaches (k-Anonymity)")
	rootCmd.AddCommand(auditCmd)
}
