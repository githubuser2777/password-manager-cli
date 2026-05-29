package crypto

import (
	"bytes"
	"strings"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	salt, _ := GenerateSalt(16)
	pw := []byte("my-master-password")
	key1 := DeriveKey(pw, salt)
	key2 := DeriveKey(pw, salt)

	if len(key1) != 32 {
		t.Fatalf("expected key length 32, got %d", len(key1))
	}
	if !bytes.Equal(key1, key2) {
		t.Fatal("keys derived from same password and salt should be identical")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key, _ := GenerateSalt(32) // using GenerateSalt to get 32 bytes of random key
	nonce, _ := GenerateNonce()
	plaintext := []byte("secret vault data")

	ciphertext, err := Encrypt(plaintext, key, nonce)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key, nonce)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestGenerateRandomPassword(t *testing.T) {
	pw, err := GenerateRandomPassword(16, true)
	if err != nil {
		t.Fatalf("failed to generate password: %v", err)
	}
	if len(pw) != 16 {
		t.Fatalf("expected password length 16, got %d", len(pw))
	}

	// Verify guaranteed character distribution
	hasUpper := strings.ContainsAny(pw, upperChars)
	hasLower := strings.ContainsAny(pw, lowerChars)
	hasDigit := strings.ContainsAny(pw, numberChars)
	hasSpecial := strings.ContainsAny(pw, specialChars)

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		t.Fatalf("generated password %q does not satisfy guaranteed classes", pw)
	}
}

func TestValidateMasterPassword(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{"Short1!", false},
		{"NoSpecialChar123", false},
		{"nospecialchar123!", false},
		{"NOSPECIALCHAR123!", false},
		{"OnlyLettersNoDigits!", false},
		{"ValidMasterPassword123!", true},
	}

	for _, tt := range tests {
		err := ValidateMasterPassword([]byte(tt.password))
		if tt.valid && err != nil {
			t.Errorf("expected password %q to be valid, got error: %v", tt.password, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("expected password %q to be invalid, but got no error", tt.password)
		}
	}
}
