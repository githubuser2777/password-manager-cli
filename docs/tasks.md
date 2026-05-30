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

## Phase 6: Comprehensive Security & TUI Revamp Enhancements
- [x] **Task 6.1 (Cryptography Module):**
  - [x] Implement `ZeroBytes(b []byte)` in `internal/crypto/crypto.go`.
  - [x] Refactor `DeriveKey` to accept `[]byte` instead of `string`.
  - [x] Implement `ValidateMasterPassword(pw []byte) error` complexity checks.
  - [x] Rewrite `GenerateRandomPassword` to guarantee character classes and implement cryptographically secure shuffle.
- [x] **Task 6.2 (Storage Module - Dynamic V2 Headers & Compatibility):**
  - [x] Update `SaveVault` to output the V2 binary structure with `"PMV2"` magic signature and serialized Argon2id params.
  - [x] Update `LoadVault` to recognize V2 magic header and fall back to legacy V1 formats automatically.
  - [x] Ensure all key and plain JSON buffers are securely zeroed out after use.
- [x] **Task 6.3 (CLI Improvements):**
  - [x] Refactor `promptPassword` in `cmd/utils.go` to return `[]byte` and zero out raw strings.
  - [x] Integrate strength validation into `cmd/init.go` and `cmd/changepass.go` with strict memory zeroing.
  - [x] Update `cmd/get.go` with non-blocking interactive console countdown for clipboard clearing, including OS interrupt handling.
- [x] **Task 6.4 (Premium TUI Revamp & Interactive Flow):**
  - [x] Implement the `stateDecrypting` transition and incorporate the Purple/Lavender `spinner` loader bubble.
  - [x] Setup background asynchronous goroutine job for key derivation & vault decrypt.
  - [x] Implement split-screen dashboard layout with a detailed card panel on the right (and responsive width check).
  - [x] Create real-time visual password strength progress bar (Weak/Medium/Strong with color coding) in form states.
  - [x] Implement asynchronous TUI clipboard countdown timer using `tea.Tick` and render status footer.
- [x] **Task 6.5 (Verification & Testing):**
  - [x] Update unit tests for `crypto` and `storage` to achieve coverage for new validations, dynamic params, and legacy vault migrations.
  - [x] Compile and verify binary functionality, manual testing of revamped TUI dashboard, spinner, and clipboard auto-clears.

## Phase 7: Refactoring & Code Quality
- [x] **Task 7.1 (Testing Strategy - core & cmd):**
  - [x] Write unit tests for `internal/core` struct integrity.
  - [x] Refactor `cmd` functions to separate cobra boilerplate from testable business logic.
  - [x] Write unit tests for the core CLI commands using mock stdin/stdout and temporary directories.
- [x] **Task 7.2 (Testing Strategy - tui):**
  - [x] Write Bubble Tea update loop tests for `internal/tui`.
- [x] **Task 7.3 (Linting & Optimizations):**
  - [x] Run `go fmt ./...` and `go vet ./...` natively and fix any issues found.
  - [x] Review data structures for memory leakage or unnecessary copies of sensitive fields.
