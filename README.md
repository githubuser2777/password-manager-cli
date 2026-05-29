# Password Manager CLI

A fast, secure, and lightweight command-line password manager written in Go.

## Features
- **Interactive TUI**: A beautiful, modern Terminal User Interface (powered by Bubble Tea) to easily view, add, edit, delete, and search your passwords.
- **AES-256-GCM Encryption**: State-of-the-art authenticated encryption.
- **Argon2id Key Derivation**: High security against brute-force attacks on your Master Password.
- **Local Storage**: Your vault is stored entirely on your local machine (default: `~/.passmgr/vault.enc`). No cloud dependencies.
- **Clipboard Integration**: Seamlessly copy passwords to your clipboard without displaying them on the screen.
- **Random Password Generator**: Built-in utility to generate strong, secure random passwords on the fly.

## Installation

Ensure you have [Go](https://go.dev/dl/) installed. 
Open your terminal and run the following command to build from source:
```bash
go build -o passmgr.exe
```
*(On Linux/macOS, use `go build -o passmgr`)*

## Usage Guide

### 1. Interactive TUI (Recommended)
The easiest way to use the password manager is via the built-in Terminal User Interface. Simply run the executable without any arguments:

```bash
passmgr
```
*(Or `.\passmgr.exe` on Windows)*

- **Unlock**: Enter your master password to unlock the vault.
- **List & Search**: Browse your saved services or press `/` to search.
- **View & Copy**: Press `Enter` on a service to view details (including Username, Password, Notes, and Timestamps). Press `c` to securely copy the password to your clipboard.
- **Add**: Press `a` from the list to open the interactive Add form.
- **Edit**: Press `e` on a selected service to edit its details.
- **Delete**: Press `d` to delete a service (prompts for confirmation).
- **Audit Passwords**: Press `r` to run a local security check (finds short, reused passwords). Press `R` (Shift+R) to perform a deep online audit against known data breaches via HaveIBeenPwned API.

### 2. CLI Commands (Scripting & Automation)
If you prefer traditional commands or want to write automated scripts, all features are available via the CLI:

- **Initialize Vault**: `passmgr init`
- **Add Credential**: `passmgr add github.com [-g]` (use `-g` to auto-generate password)
- **Get Password**: `passmgr get github.com [-c]` (use `-c` to copy to clipboard)
- **Update Credential**: `passmgr update github.com`
- **Delete Credential**: `passmgr delete github.com`
- **Search Services**: `passmgr search "git"`
- **List All**: `passmgr list`
- **Generate Password**: `passmgr generate 20`
- **Export/Import Vault**: `passmgr export`, `passmgr import`
- **Audit Passwords**: `passmgr audit`
- **Change Master Password**: `passmgr changepass`

## Build Scripts
You can use the provided automated scripts to build the application:
- **Windows**: Run `.\build.ps1`
- **Linux/macOS**: Run `make build` (Requires `make` to be installed)
