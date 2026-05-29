# Execution Checklist

Rules: 
1. Upon completing a task, the Code Agent will mark it with an `[x]`. If a task is too complex, proactively break it down further in this file.
2. Strictly adhere to `rules.md`: **Do not open/read more than 5 files at once**; **Mandatory comments** for complex cryptographic logic.

## Phase 1: Foundation
- [x] Initialize directory and Go module (`go mod init password-manager-cli`).
- [x] Set up directory structure: `cmd/`, `internal/crypto/`, `internal/storage/`, `internal/core/`.
- [x] Install the `spf13/cobra` package and set up `main.go`, `cmd/root.go`.

## Phase 2: Domain Models & Cryptography Module
- [x] (internal/core): Define `Vault` and `Entry` structs.
- [x] (internal/crypto): Implement `DeriveKey(masterPassword string, salt []byte) []byte` using Argon2id. Implement `GenerateSalt() []byte`.
- [x] (internal/crypto): Implement `Encrypt(plaintext, key, nonce []byte) ([]byte, error)` and `Decrypt(ciphertext, key, nonce []byte) ([]byte, error)` using AES-256-GCM.
- [x] (internal/crypto): Implement `GenerateRandomPassword(length int, includeSpecial bool) string`.
- [x] Write Unit Tests to verify the accuracy of the `crypto` module (Encrypting and decrypting must yield the original data).

## Phase 3: Storage Module
- [x] (internal/storage): Build the `SaveVault(path string, masterPassword string, vault *core.Vault) error` function based on the Storage Flow (write binary Salt + Nonce + Ciphertext).
- [x] (internal/storage): Build the `LoadVault(path string, masterPassword string) (*core.Vault, error)` function.
- [x] Write Unit Tests for saving and loading the vault from a file (using a temporary directory).

## Phase 4: CLI Commands
- [x] (cmd/init.go): Implement `init` command. (Check if vault exists, prompt for hidden Master Password twice to confirm, create an empty vault with a new Salt and save it).
- [x] (cmd/add.go): Implement `add <service>` command. (Input username, input password or auto-generate if using the `--generate` flag).
- [x] (cmd/get.go): Implement `get <service>` command. (Display username and password. Support the `--copy` flag to auto-copy to clipboard instead of printing to screen).
- [x] (cmd/list.go): Implement `list` command. (Print all saved services).
- [x] (cmd/generate.go): Implement `generate` command. (Generate a random password for immediate use).

## Phase 5: Polish & Build
- [x] Package and write build scripts (`Makefile` or `build.ps1`).
- [x] Write the `README.md` document with instructions for using the tool.
