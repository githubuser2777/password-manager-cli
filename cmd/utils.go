package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/term"
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

	start := 0
	end := len(bytePassword)
	for start < end && isSpace(bytePassword[start]) {
		start++
	}
	for end > start && isSpace(bytePassword[end-1]) {
		end--
	}

	trimmed := make([]byte, end-start)
	copy(trimmed, bytePassword[start:end])

	// Zero out original raw read buffer
	for i := range bytePassword {
		bytePassword[i] = 0
	}

	return trimmed, nil
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\v' || b == '\f' || b == '\r'
}
