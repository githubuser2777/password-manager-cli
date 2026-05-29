package storage

import (
	"encoding/json"
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
	masterPw := []byte("super_secret_master_password_123!")

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

	// Verify V2 signature
	data, err := os.ReadFile(vaultPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if len(data) < 4 || string(data[:4]) != "PMV2" {
		t.Errorf("expected modern vault to start with magic 'PMV2'")
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
	_, err = LoadVault(vaultPath, []byte("wrong_password"))
	if err == nil {
		t.Fatal("LoadVault should have failed with wrong password")
	}
}

func TestLegacyV1VaultCompatibility(t *testing.T) {
	// Let's manually construct a V1 legacy vault file
	// Layout V1: Salt(16) + Nonce(12) + Ciphertext
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "legacy_vault.enc")
	masterPw := []byte("legacy_pass_123")

	salt, _ := crypto.GenerateSalt(16)
	nonce, _ := crypto.GenerateNonce()
	
	// Encode JSON data
	dummyVault := core.Vault{
		Salt: salt,
		Entries: map[string]core.Entry{
			"legacy.com": {
				Username: "legacy_user",
				Password: "legacy_password",
			},
		},
	}
	jsonData, _ := json.Marshal(&dummyVault)

	// Derive V1 key using standard defaults
	key := crypto.DeriveKeyWithParams(masterPw, salt, 1, 64*1024, 4)
	ciphertext, _ := crypto.Encrypt(jsonData, key, nonce)

	// Combine to build V1 file data
	var fileData []byte
	fileData = append(fileData, salt...)
	fileData = append(fileData, nonce...)
	fileData = append(fileData, ciphertext...)

	err := os.WriteFile(vaultPath, fileData, 0600)
	if err != nil {
		t.Fatalf("failed to write legacy vault file: %v", err)
	}

	// Now try to load this legacy vault using our LoadVault which supports auto-fallback
	loadedVault, err := LoadVault(vaultPath, masterPw)
	if err != nil {
		t.Fatalf("LoadVault failed to read legacy V1 vault: %v", err)
	}

	if loadedVault.Entries["legacy.com"].Password != "legacy_password" {
		t.Errorf("expected legacy_password, got %s", loadedVault.Entries["legacy.com"].Password)
	}

	// Upgrade Legacy Vault to V2
	err = SaveVault(vaultPath, masterPw, loadedVault)
	if err != nil {
		t.Fatalf("failed to upgrade legacy vault to V2: %v", err)
	}

	// Check that it now starts with PMV2
	newData, _ := os.ReadFile(vaultPath)
	if len(newData) < 4 || string(newData[:4]) != "PMV2" {
		t.Errorf("expected upgraded vault to start with magic 'PMV2'")
	}

	// Try reading again
	rereadVault, err := LoadVault(vaultPath, masterPw)
	if err != nil {
		t.Fatalf("failed to read upgraded V2 vault: %v", err)
	}
	if rereadVault.Entries["legacy.com"].Password != "legacy_password" {
		t.Errorf("expected legacy_password after upgrade, got %s", rereadVault.Entries["legacy.com"].Password)
	}
}
