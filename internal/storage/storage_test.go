package storage

import (
	"os"
	"path/filepath"
	"testing"

	"password-manager-cli/internal/core"
	"password-manager-cli/internal/crypto"
)

func TestSaveAndLoadVault(t *testing.T) {
	// Setup temporary file
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "vault.enc")
	masterPw := "super_secret_master_password"

	// Create a dummy vault
	salt, _ := crypto.GenerateSalt(16)
	originalVault := &core.Vault{
		Salt: salt,
		Entries: map[string]core.Entry{
			"github.com": {
				Username:  "admin",
				Password:  "password123",
				CreatedAt: "2026-05-29",
			},
		},
	}

	// Save
	err := SaveVault(vaultPath, masterPw, originalVault)
	if err != nil {
		t.Fatalf("SaveVault failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		t.Fatalf("Vault file was not created")
	}

	// Load successfully
	loadedVault, err := LoadVault(vaultPath, masterPw)
	if err != nil {
		t.Fatalf("LoadVault failed: %v", err)
	}

	if loadedVault.Entries["github.com"].Password != "password123" {
		t.Errorf("expected password123, got %s", loadedVault.Entries["github.com"].Password)
	}

	// Load with wrong password
	_, err = LoadVault(vaultPath, "wrong_password")
	if err == nil {
		t.Fatal("LoadVault should have failed with wrong password")
	}
}
