package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/term"
	"password-manager-cli/internal/vault"
)

// getVaultPath returns the absolute path to the vault.enc file.
// Defaults to ~/.passmgr/vault.enc
var getVaultPath = func() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting home directory:", err)
		os.Exit(1)
	}
	return filepath.Join(home, ".passmgr", "vault.enc")
}

// promptPassword securely prompts the user for a password without echoing.
// Returns a byte slice and zeroes out the raw temporary buffer.
var promptPassword = func(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // print a newline after reading
	if err != nil {
		return nil, err
	}

	trimmed := bytes.Clone(bytes.TrimSpace(bytePassword))

	// Zero out original raw read buffer
	clear(bytePassword)

	return trimmed, nil
}

func withVault(action func(v *vault.Vault, masterPw []byte, path string)) {
	path := getVaultPath()
	masterPw, err := promptPassword("Master Password: ")
	if err != nil {
		return
	}
	defer clear(masterPw)

	v, err := vault.LoadVault(path, masterPw)
	if err != nil {
		fmt.Println("Failed to open vault:", err)
		return
	}
	action(v, masterPw, path)
}

func withVaultMutate(action func(v *vault.Vault, masterPw []byte, path string) bool) {
	withVault(func(v *vault.Vault, masterPw []byte, path string) {
		if action(v, masterPw, path) {
			if err := vault.SaveVault(path, masterPw, v); err != nil {
				fmt.Println("Failed to save vault:", err)
			}
		}
	})
}
