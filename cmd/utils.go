package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// getVaultPath returns the absolute path to the vault.enc file.
// Defaults to ~/.passmgr/vault.enc
func getVaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting home directory:", err)
		os.Exit(1)
	}
	return filepath.Join(home, ".passmgr", "vault.enc")
}

// promptPassword securely prompts the user for a password without echoing.
func promptPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // print a newline after reading
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePassword)), nil
}
