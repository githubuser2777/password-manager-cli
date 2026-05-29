# Product Specification: Password Manager CLI

## 1. Overview
A command-line tool (CLI) built with Go that helps generate highly secure random passwords and store them safely. The entire password vault is protected and decrypted via a single Master Password.

## 2. Target Audience
- Individuals, system administrators, and developers who need a fast, secure password manager that operates entirely within the terminal.

## 3. Core Features
- `init`: Initialize the initial vault storage and set the Master Password.
- `add`: Add an account (Website/App, Username, Password). Option to auto-generate a password for this account.
- `get`: Retrieve and decrypt a password based on the name (requires Master Password authentication). Supports copying directly to the clipboard.
- `list`: View a list of saved services (does not display actual passwords).
- `generate`: Generate a secure random password, allowing customization of length and character types.

## 4. Security Requirements
- **Encryption**: Data must be encrypted using AES-256 in GCM mode (ensuring both data confidentiality and integrity).
- **Key Derivation**: Use the `Argon2id` algorithm to hash the Master Password, defending against Brute-force and Dictionary attacks.
- **Storage**: Passwords are encrypted as a monolith and stored locally on the user's computer (e.g., `~/.passmgr/vault.enc`).
- **Memory Security**: When entering the Master Password in the terminal, characters must not be displayed on the screen.
