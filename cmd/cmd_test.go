package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCmd(t *testing.T) {
	tmpDir := t.TempDir()
	vaultFile := filepath.Join(tmpDir, "vault.enc")

	oldVaultPath := getVaultPath
	getVaultPath = func() string { return vaultFile }
	defer func() { getVaultPath = oldVaultPath }()

	oldPrompt := promptPassword
	promptIdx := 0
	prompts := [][]byte{
		[]byte("StrongPassw0rd!"),
		[]byte("StrongPassw0rd!"),
	}
	promptPassword = func(prompt string) ([]byte, error) {
		p := prompts[promptIdx]
		promptIdx++
		return p, nil
	}
	defer func() { promptPassword = oldPrompt }()

	initCmd.Run(initCmd, []string{})

	if _, err := os.Stat(vaultFile); os.IsNotExist(err) {
		t.Errorf("Expected vault file to be created at %s", vaultFile)
	}

	// Test initializing again (should just return)
	initCmd.Run(initCmd, []string{})
}

func TestAddCmd(t *testing.T) {
	tmpDir := t.TempDir()
	vaultFile := filepath.Join(tmpDir, "vault.enc")

	oldVaultPath := getVaultPath
	getVaultPath = func() string { return vaultFile }
	defer func() { getVaultPath = oldVaultPath }()

	oldPrompt := promptPassword
	promptPassword = func(prompt string) ([]byte, error) {
		if prompt == "Master Password: " {
			return []byte("StrongPassw0rd!"), nil
		}
		if prompt == "Password: " {
			return []byte("ServicePass123!"), nil
		}
		return []byte("StrongPassw0rd!"), nil
	}
	defer func() { promptPassword = oldPrompt }()

	// Initialize first
	initCmd.Run(initCmd, []string{})

	// Mock stdin for addCmd
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	go func() {
		w.Write([]byte("testuser\ntestnotes\n"))
		w.Close()
	}()

	generateFlag = false
	addCmd.Run(addCmd, []string{"github.com"})

	// Verify we can't add it again (returns early)
	addCmd.Run(addCmd, []string{"github.com"})
}

func TestGetCmd(t *testing.T) {
	tmpDir := t.TempDir()
	vaultFile := filepath.Join(tmpDir, "vault.enc")

	oldVaultPath := getVaultPath
	getVaultPath = func() string { return vaultFile }
	defer func() { getVaultPath = oldVaultPath }()

	oldPrompt := promptPassword
	promptPassword = func(prompt string) ([]byte, error) {
		return []byte("StrongPassw0rd!"), nil
	}
	defer func() { promptPassword = oldPrompt }()

	initCmd.Run(initCmd, []string{})

	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	go func() {
		w.Write([]byte("testuser\ntestnotes\n"))
		w.Close()
	}()

	generateFlag = false
	addCmd.Run(addCmd, []string{"github.com"})

	// Now run get
	copyFlag = false
	getCmd.Run(getCmd, []string{"github.com"})

	// And get with a non-existent service
	getCmd.Run(getCmd, []string{"nonexistent.com"})
}

func TestOtherCmds(t *testing.T) {
	tmpDir := t.TempDir()
	vaultFile := filepath.Join(tmpDir, "vault.enc")

	oldVaultPath := getVaultPath
	getVaultPath = func() string { return vaultFile }
	defer func() { getVaultPath = oldVaultPath }()

	oldPrompt := promptPassword
	promptPassword = func(prompt string) ([]byte, error) {
		return []byte("StrongPassw0rd!"), nil
	}
	defer func() { promptPassword = oldPrompt }()

	initCmd.Run(initCmd, []string{})

	// Test generate
	generateCmd.Run(generateCmd, []string{"16"})

	// Add a service
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	go func() {
		w.Write([]byte("username\nnotes\n"))
	}()
	generateFlag = false
	addCmd.Run(addCmd, []string{"github.com"})

	// Test list
	listCmd.Run(listCmd, []string{})

	// Test search
	searchCmd.Run(searchCmd, []string{"git"})

	// Test delete
	go func() {
		w.Write([]byte("y\n"))
		w.Close()
	}()
	deleteCmd.Run(deleteCmd, []string{"github.com"})

	os.Stdin = oldStdin
}
