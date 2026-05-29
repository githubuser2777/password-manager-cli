package crypto

import (
	"bytes"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	salt, _ := GenerateSalt(16)
	key1 := DeriveKey("my-master-password", salt)
	key2 := DeriveKey("my-master-password", salt)

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
}
