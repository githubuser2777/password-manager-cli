# Changelog

All major project changes will be recorded in this file.

## [Phase 5] - 2026-05-29
### Added
- Added the `LICENSE` file (MIT License) to open-source the project.
- Moved and rewrote `SECURITY.md` to `.github/SECURITY.md` according to standard GitHub Security Policies.
- Added `README.md` with complete instructions on how to use the tool.
- Created build scripts for Windows (`build.ps1`) and Linux/macOS (`Makefile`).
- Finalized project version 1.0.

## [Phase 4] - 2026-05-29
### Added
- Integrated the `golang.org/x/term` package for secure Master Password input, hiding it from the screen.
- Integrated the `github.com/atotto/clipboard` package for quick password copying.
- Completed the 5 main commands using Cobra:
  - `passmgr init`: Initialize a new vault with a master password.
  - `passmgr add <service>`: Add a new account, supporting the `-g` flag to generate a random password.
  - `passmgr get <service>`: Retrieve a password, supporting the `-c` flag to copy it directly to the Clipboard.
  - `passmgr list`: View a list of all saved services.
  - `passmgr generate`: Utility to quickly generate a password from the terminal.

## [Phase 3] - 2026-05-29
### Added
- Added the `internal/storage/storage.go` module to manage reading/writing IO for the `vault.enc` file.
- Built a secure Load/Save mechanism (Atomic Write): Write to a temporary file before renaming to prevent file corruption in case of a crash midway.
- Packaged binary data with the structure `Salt(16) + Nonce(12) + Ciphertext`.
- Wrote Unit Tests (`storage_test.go`) simulating read/write flows using temporary directories.

## [Phase 2] - 2026-05-29
### Added
- Added `internal/core/models.go` defining the `Vault` and `Entry` structs.
- Wrote the `internal/crypto/crypto.go` cryptography module with Argon2id (Key Derivation) and AES-256-GCM (Encrypt/Decrypt).
- Added the `GenerateRandomPassword` function supporting random password generation with optional special characters.
- Wrote Unit Tests (`crypto_test.go`) for the entire cryptography module (100% Core Logic Coverage).

## [Phase 1] - 2026-05-29
### Added
- Initialized directory and Go module (`password-manager-cli`).
- Set up standard directory structure: `cmd/`, `internal/crypto/`, `internal/storage/`, `internal/core/`.
- Installed the `spf13/cobra` package.
- Set up the `main.go` file and the root command `cmd/root.go` with basic descriptions.
- Successfully integrated GitHub Spec Kit (`.specify/`, `.github/`).

### Changed
- Updated `rules.md` to add rules regarding the Git process and Changelog.
