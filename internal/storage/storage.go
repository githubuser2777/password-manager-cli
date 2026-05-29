package storage

import (
	"encoding/binary"
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

// Dynamic Argon2id parameters that can be changed or tuned.
var (
	ArgonTime    uint32 = 1
	ArgonMemory  uint32 = 64 * 1024 // 64 MB
	ArgonThreads uint8  = 4
)

// SaveVault safely writes the vault to the specified path using the V2 format.
// It uses atomic writing (write to temp file, then rename) to prevent corruption.
func SaveVault(path string, masterPassword []byte, vault *core.Vault) error {
	if len(vault.Salt) != saltLen {
		return errors.New("invalid vault salt length")
	}

	// 1. Marshal Vault to JSON
	jsonData, err := json.Marshal(vault)
	if err != nil {
		return err
	}
	defer crypto.ZeroBytes(jsonData) // Securely zero out JSON plaintext

	// 2. Generate Nonce
	nonce, err := crypto.GenerateNonce()
	if err != nil {
		return err
	}

	// 3. Derive Key with V2 parameters
	key := crypto.DeriveKeyWithParams(masterPassword, vault.Salt, ArgonTime, ArgonMemory, ArgonThreads)
	defer crypto.ZeroBytes(key) // Securely zero out derived key

	// 4. Encrypt JSON data
	ciphertext, err := crypto.Encrypt(jsonData, key, nonce)
	if err != nil {
		return err
	}

	// 5. Construct final V2 binary format:
	// Magic "PMV2" (4 bytes) + Salt(16) + Nonce(12) + Time(4) + Memory(4) + Threads(1) + Ciphertext
	var fileData []byte
	fileData = append(fileData, []byte("PMV2")...)
	fileData = append(fileData, vault.Salt...)
	fileData = append(fileData, nonce...)

	timeBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(timeBuf, ArgonTime)
	fileData = append(fileData, timeBuf...)

	memBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(memBuf, ArgonMemory)
	fileData = append(fileData, memBuf...)

	fileData = append(fileData, ArgonThreads)
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
// It supports both V1 (Legacy) and V2 (Modern with Magic Bytes) vault formats.
func LoadVault(path string, masterPassword []byte) (*core.Vault, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Detect if file has V2 Magic prefix "PMV2"
	isV2 := false
	if len(fileData) >= 4 && string(fileData[:4]) == "PMV2" {
		isV2 = true
	}

	var salt, nonce, ciphertext []byte
	var time, memory uint32
	var threads uint8

	if isV2 {
		// Minimum size for V2 header + payload
		if len(fileData) < 4+saltLen+nonceLen+4+4+1 {
			return nil, errors.New("file is too short or corrupted")
		}

		salt = fileData[4 : 4+saltLen]
		nonce = fileData[4+saltLen : 4+saltLen+nonceLen]

		paramStart := 4 + saltLen + nonceLen
		time = binary.BigEndian.Uint32(fileData[paramStart : paramStart+4])
		memory = binary.BigEndian.Uint32(fileData[paramStart+4 : paramStart+8])
		threads = fileData[paramStart+8]

		ciphertext = fileData[paramStart+9:]
	} else {
		// Fall back to legacy V1 layout
		if len(fileData) < saltLen+nonceLen {
			return nil, errors.New("file is too short or corrupted")
		}
		salt = fileData[:saltLen]
		nonce = fileData[saltLen : saltLen+nonceLen]
		ciphertext = fileData[saltLen+nonceLen:]

		// Legacy defaults
		time = 1
		memory = 64 * 1024
		threads = 4
	}

	// 2. Derive Key
	key := crypto.DeriveKeyWithParams(masterPassword, salt, time, memory, threads)
	defer crypto.ZeroBytes(key) // Securely zero out derived key

	// 3. Decrypt
	jsonData, err := crypto.Decrypt(ciphertext, key, nonce)
	if err != nil {
		return nil, errors.New("invalid master password or corrupted data")
	}
	defer crypto.ZeroBytes(jsonData) // Securely zero out JSON plaintext

	// 4. Unmarshal
	var vault core.Vault
	if err := json.Unmarshal(jsonData, &vault); err != nil {
		return nil, err
	}

	return &vault, nil
}
