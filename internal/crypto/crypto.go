package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"math/big"

	"golang.org/x/crypto/argon2"
)

// GenerateSalt creates a secure random salt of the given length.
// A 16-byte salt is standard for Argon2id.
func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// DeriveKey generates a 256-bit (32 bytes) key from a master password and salt using Argon2id.
// It uses recommended parameters: time=1, memory=64MB, threads=4.
func DeriveKey(masterPassword string, salt []byte) []byte {
	time := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLen := uint32(32)

	return argon2.IDKey([]byte(masterPassword), salt, time, memory, threads, keyLen)
}

// GenerateNonce creates a 12-byte random nonce for AES-GCM.
func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

// Encrypt encrypts plaintext using AES-256-GCM.
// The key must be exactly 32 bytes long. The nonce must be exactly 12 bytes long.
func Encrypt(plaintext, key, nonce []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes for AES-256")
	}
	if len(nonce) != 12 {
		return nil, errors.New("nonce must be 12 bytes for AES-GCM")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Seal appends the encrypted data and auth tag to the first argument (dst)
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts AES-256-GCM ciphertext.
// The key must be exactly 32 bytes long. The nonce must be exactly 12 bytes long.
func Decrypt(ciphertext, key, nonce []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes for AES-256")
	}
	if len(nonce) != 12 {
		return nil, errors.New("nonce must be 12 bytes for AES-GCM")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars  = "0123456789"
	specialChars = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

// GenerateRandomPassword creates a secure random password.
func GenerateRandomPassword(length int, includeSpecial bool) (string, error) {
	if length < 4 {
		return "", errors.New("password length must be at least 4")
	}

	charSet := lowerChars + upperChars + numberChars
	if includeSpecial {
		charSet += specialChars
	}

	password := make([]byte, length)
	charSetLen := big.NewInt(int64(len(charSet)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charSetLen)
		if err != nil {
			return "", err
		}
		password[i] = charSet[num.Int64()]
	}

	return string(password), nil
}
