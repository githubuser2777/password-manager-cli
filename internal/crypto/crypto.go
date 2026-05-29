package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"math/big"
	"strings"

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

// DeriveKey generates a 256-bit (32 bytes) key from a master password and salt using Argon2id with default parameters.
func DeriveKey(masterPassword []byte, salt []byte) []byte {
	return DeriveKeyWithParams(masterPassword, salt, 1, 64*1024, 4)
}

// DeriveKeyWithParams generates a 256-bit key using custom Argon2id parameters.
func DeriveKeyWithParams(masterPassword []byte, salt []byte, time, memory uint32, threads uint8) []byte {
	return argon2.IDKey(masterPassword, salt, time, memory, threads, 32)
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

// ZeroBytes overwrites the byte slice with zeros to clear sensitive data from memory.
func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// ValidateMasterPassword verifies that a master password meets minimum security guidelines:
// - At least 12 characters
// - At least one uppercase letter
// - At least one lowercase letter
// - At least one digit
// - At least one special character
func ValidateMasterPassword(pw []byte) error {
	if len(pw) < 12 {
		return errors.New("master password must be at least 12 characters long")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range string(pw) { // Cast to string to safely parse runes
		switch {
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	if !hasLower {
		return errors.New("master password must contain at least one lowercase letter")
	}
	if !hasUpper {
		return errors.New("master password must contain at least one uppercase letter")
	}
	if !hasDigit {
		return errors.New("master password must contain at least one digit")
	}
	if !hasSpecial {
		return errors.New("master password must contain at least one special character")
	}

	return nil
}

// GenerateRandomPassword creates a secure random password.
// It guarantees that at least one character from each active character set is included
// and uses a secure shuffle to prevent pattern predictability.
func GenerateRandomPassword(length int, includeSpecial bool) (string, error) {
	if length < 4 {
		return "", errors.New("password length must be at least 4")
	}

	// 1. Prepare required characters
	var required []byte
	required = append(required, lowerChars[mustRandInt(len(lowerChars))])
	required = append(required, upperChars[mustRandInt(len(upperChars))])
	required = append(required, numberChars[mustRandInt(len(numberChars))])
	if includeSpecial {
		required = append(required, specialChars[mustRandInt(len(specialChars))])
	}

	if len(required) > length {
		return "", errors.New("password length is too short for the required character classes")
	}

	// 2. Generate the remaining characters
	charSet := lowerChars + upperChars + numberChars
	if includeSpecial {
		charSet += specialChars
	}

	password := make([]byte, length)
	copy(password, required)

	charSetLen := big.NewInt(int64(len(charSet)))
	for i := len(required); i < length; i++ {
		num, err := rand.Int(rand.Reader, charSetLen)
		if err != nil {
			return "", err
		}
		password[i] = charSet[num.Int64()]
	}

	// 3. Cryptographically secure Fisher-Yates shuffle
	for i := length - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		j := int(jBig.Int64())
		password[i], password[j] = password[j], password[i]
	}

	res := string(password)
	ZeroBytes(password) // Zero the temporary buffer
	return res, nil
}

// mustRandInt is a helper to securely select a random index, panicking on crypto error
func mustRandInt(max int) int {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic("cryptographic failure: " + err.Error())
	}
	return int(val.Int64())
}
