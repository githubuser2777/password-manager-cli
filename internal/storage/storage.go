package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"password-manager-cli/internal/core"
	"password-manager-cli/internal/crypto"
)

const (
	saltLen  = 16
	nonceLen = 12
)

// SaveVault safely writes the vault to the specified path.
// It uses atomic writing (write to temp file, then rename) to prevent corruption.
func SaveVault(path string, masterPassword string, vault *core.Vault) error {
	if len(vault.Salt) != saltLen {
		return errors.New("invalid vault salt length")
	}

	// 1. Marshal Vault to JSON
	jsonData, err := json.Marshal(vault)
	if err != nil {
		return err
	}

	// 2. Generate Nonce
	nonce, err := crypto.GenerateNonce()
	if err != nil {
		return err
	}

	// 3. Derive Key
	key := crypto.DeriveKey(masterPassword, vault.Salt)

	// 4. Encrypt JSON data
	ciphertext, err := crypto.Encrypt(jsonData, key, nonce)
	if err != nil {
		return err
	}

	// 5. Construct final binary format: Salt(16) + Nonce(12) + Ciphertext
	var fileData []byte
	fileData = append(fileData, vault.Salt...)
	fileData = append(fileData, nonce...)
	fileData = append(fileData, ciphertext...)

	// 6. Write securely (Atomic Write)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	tempFile := path + ".tmp"
	if err := os.WriteFile(tempFile, fileData, 0600); err != nil {
		return err
	}

	// Create backup of existing vault if it exists
	if _, err := os.Stat(path); err == nil {
		backupPath := path + ".bak"
		if input, err := os.ReadFile(path); err == nil {
			_ = os.WriteFile(backupPath, input, 0600)
		}
	}

	// Rename over the original file
	return os.Rename(tempFile, path)
}

// LoadVault reads and decrypts the vault from the specified path.
func LoadVault(path string, masterPassword string) (*core.Vault, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(fileData) < saltLen+nonceLen {
		return nil, errors.New("file is too short or corrupted")
	}

	// 1. Extract Salt, Nonce, Ciphertext
	salt := fileData[:saltLen]
	nonce := fileData[saltLen : saltLen+nonceLen]
	ciphertext := fileData[saltLen+nonceLen:]

	// 2. Derive Key
	key := crypto.DeriveKey(masterPassword, salt)

	// 3. Decrypt
	jsonData, err := crypto.Decrypt(ciphertext, key, nonce)
	if err != nil {
		return nil, errors.New("invalid master password or corrupted data")
	}

	// 4. Unmarshal
	var vault core.Vault
	if err := json.Unmarshal(jsonData, &vault); err != nil {
		return nil, err
	}

	return &vault, nil
}
