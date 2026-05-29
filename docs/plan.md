# Technical Blueprint: Password Manager CLI

## 1. Tech Stack
- **Language**: Go (Golang)
- **CLI Framework**: `github.com/spf13/cobra` for command structure and `golang.org/x/term` to hide passwords during input.
- **Cryptography**: 
  - Key Derivation: `golang.org/x/crypto/argon2`
  - Encryption/Decryption: `crypto/aes` and `crypto/cipher` (Go standard libraries).
  - Secure Random Number Generation: `crypto/rand`
- **Storage**: Local file, content is a JSON structure fully encrypted with AES-256.
- **Utilities (Optional)**: `github.com/atotto/clipboard` to copy passwords to the clipboard.

## 2. Architecture & Modules
The project applies the `Standard Go Project Layout`:
- `cmd/`: Contains definitions and entrypoints for CLI commands (`root.go`, `init.go`, `add.go`, `get.go`, `list.go`, `generate.go`).
- `internal/crypto/`: Encapsulates all cryptographic logic (AES-GCM, Argon2id, Random Password).
- `internal/storage/`: Manages IO (reading/writing local files) and JSON operations.
- `internal/core/`: Contains the core domain models (Entities) of the application.

## 3. Data Models

### Entry (A single account)
```go
type Entry struct {
    Username  string `json:"username"`
    Password  string `json:"password"`
    CreatedAt string `json:"created_at"` // ISO8601 string
}
```

### Vault (The main storage)
```go
type Vault struct {
    // Salt shared for Argon2 hashing (necessary to regenerate the Key). Will be randomly generated on init.
    Salt    []byte           `json:"salt"` 
    Entries map[string]Entry `json:"entries"` // Key is the service name (e.g., "github", "google")
}
```

## 4. Storage Flow
- **`vault.enc` File Structure**:
  Since `Salt` must be known beforehand to hash the Master Password, we have two directions:
  - Direction 1: `vault.enc` contains the `Salt` at the beginning of the file unencrypted (plaintext), followed by the `Nonce` and the Encrypted Data (Ciphertext).
  - Direction 2: Store metadata (Salt, Nonce) and Ciphertext as a JSON wrapper outside.
  => **Chosen Direction 1** (Binary format) for optimization and security (Salt (16 bytes) + Nonce (12 bytes) + Ciphertext).

- **Read (Load)**:
  1. Read the `vault.enc` file.
  2. Extract the first 16 bytes as `Salt`, the next 12 bytes as `Nonce`, and the rest as `Ciphertext`.
  3. Pass `Salt` + Master Password through Argon2id -> `256-bit Key`.
  4. Decrypt `Ciphertext` with `Key` and `Nonce` using AES-GCM -> JSON string.
  5. JSON Unmarshal into the `Vault` struct.

- **Write (Save)**:
  1. JSON Marshal the `Vault` struct.
  2. Generate a new `Nonce` (12 bytes).
  3. Encrypt the JSON string using AES-GCM -> `Ciphertext`.
  4. Write to file: `Salt` + `Nonce` + `Ciphertext`.

## 5. Coding Standards & Conventions (According to rules.md)
- **Context Management**: The Code Agent must not read more than 5 files at once.
- **Code Standard**: Write clear code, prioritizing readability over overly concise code.
- **Comments**: Comments explaining complex business logic are strictly required, especially in the `internal/crypto` module (handling Argon2id, AES-GCM).
