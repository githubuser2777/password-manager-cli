package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"password-manager-cli/internal/crypto"
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
		defer crypto.ZeroBytes(masterPw)

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
		if entry.Notes != "" {
			fmt.Println("Notes:", entry.Notes)
		}
		if entry.UpdatedAt != "" {
			fmt.Println("Last Updated:", entry.UpdatedAt)
		}

		if copyFlag {
			if err := clipboard.WriteAll(entry.Password); err != nil {
				fmt.Println("Failed to copy to clipboard:", err)
			} else {
				// 1. Setup Signal Interrupt Handling
				sigChan := make(chan os.Signal, 1)
				signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
				defer signal.Stop(sigChan)

				totalDuration := 30 * time.Second
				tickInterval := 100 * time.Millisecond
				steps := int(totalDuration / tickInterval)

				fmt.Println("Password copied to clipboard! It will be cleared in 30 seconds...")

				interrupted := false
				for i := 0; i <= steps; i++ {
					select {
					case <-sigChan:
						interrupted = true
						break
					case <-time.After(tickInterval):
						remaining := float64(steps-i) * tickInterval.Seconds()
						if remaining < 0 {
							remaining = 0
						}
						percent := float64(i) / float64(steps)
						barWidth := 30
						filledWidth := int(percent * float64(barWidth))
						if filledWidth > barWidth {
							filledWidth = barWidth
						}
						bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", barWidth-filledWidth)
						fmt.Printf("\r[%s] %.1fs remaining (Ctrl+C to clear and exit)", bar, remaining)
					}
					if interrupted {
						break
					}
				}

				// 2. Overwrite Clipboard securely
				_ = clipboard.WriteAll("")
				fmt.Println("\nClipboard cleared.")
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
