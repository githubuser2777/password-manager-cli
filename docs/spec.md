# Product Specification: Password Manager CLI

## 1. Overview
A command-line tool (CLI) built with Go that helps generate highly secure random passwords and store them safely. The entire password vault is protected and decrypted via a single Master Password.

## 2. Target Audience
- Individuals, system administrators, and developers who need a fast, secure password manager that operates entirely within the terminal.

## 3. Core Features & Enhancements
- `init`: Initialize the vault storage, validate Master Password strength, and set the Master Password.
- `add`: Add an account (Website/App, Username, Password, Notes). Option to auto-generate a guaranteed-strong password for this account.
- `get`: Retrieve and decrypt a password based on the name. Supports copying to the clipboard with an asynchronous background clear mechanism in the CLI.
- `list`: View a list of saved services.
- `generate`: Generate a secure random password with a guaranteed distribution of character types (at least 1 uppercase, 1 lowercase, 1 digit, and 1 special character if enabled).
- `changepass`: Change the Master Password with strength validation.
- `audit`: Audit passwords for strength, reuse, and data breaches (via HaveIBeenPwned API).
- `tui`: A revamped, highly visual terminal user interface supporting:
  - **Modern Split-Screen Dashboard Layout**: Scrollable list of accounts on the left, beautifully formatted details card (including metadata and custom styling) on the right.
  - **Visual Password Strength Meter**: Dynamic color-coded bar displaying password strength (Weak, Medium, Strong) in real-time as the user types in the Add/Edit form.
  - **Premium Dark Palette Aesthetics**: Sleek lavender, orchid, slate-gray, and mint accents using Lipgloss styling.
  - Interactive login with visual feedback (spinner/loading) during Argon2id key derivation.
  - CRUD operations on credentials.
  - Safe clipboard copy with a non-blocking 30-second auto-clear mechanism and countdown indicator.
  - Local and online password auditing displayed in structured tables.

## 4. Security Requirements
- **Encryption**: Data encrypted using AES-256 in GCM mode.
- **Key Derivation**: Use `Argon2id` to hash the Master Password. The parameters (Time, Memory, Threads) must be dynamically encoded in the file header to allow future security upgrades and device-specific tuning without breaking backwards compatibility.
- **Memory Security**: Prompt passwords securely, handle keys/passwords in `[]byte` buffers instead of immutable `string` types where possible, and zero out memory buffers immediately after use.
- **Clipboard Security**: Any password copied to the clipboard must be securely overwritten with a blank value after exactly 30 seconds, both in the CLI and the TUI, without blocking the user interface.
- **Storage**: Passwords are encrypted and saved using an atomic write mechanism (temp file write followed by rename) to prevent corruption.

## 5. Refactoring & Code Quality Goals
- **Code Coverage**: Write comprehensive unit tests for `cmd`, `internal/core`, and `internal/tui` to bring coverage up from 0% to a healthy standard.
- **Linting & Code Formatting**: Ensure the codebase adheres to strict Go linting standards (e.g. `go vet`, `staticcheck`). 
- **Optimizations**: Perform code reviews to identify memory optimization opportunities, especially concerning sensitive data handling.
