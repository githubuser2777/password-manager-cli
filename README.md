# Password Manager CLI

A fast, secure, and lightweight command-line password manager written in Go.

## Features
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

### 1. Initialize the Vault
```bash
passmgr init
```
This command prompts you to create a Master Password. **Do not forget this password!** It is the only key to unlock your vault, and without it, your data cannot be recovered.

### 2. Add a New Credential
```bash
passmgr add github.com
```
You will be prompted for a Username and Password. Alternatively, you can let the CLI auto-generate a strong password by using the `-g` (generate) flag:
```bash
passmgr add github.com -g
```

### 3. Retrieve a Password
```bash
passmgr get github.com
```
To copy the password directly to your clipboard securely (without displaying it on the screen):
```bash
passmgr get github.com -c
```

### 4. List Saved Services
```bash
passmgr list
```
Displays a list of all services currently saved in your vault.

### 5. Generate a Random Password (Utility)
```bash
passmgr generate 20
```
*(Generates a 20-character secure password. You can omit special characters by appending the `-n` flag).*

## Build Scripts
You can use the provided automated scripts to build the application:
- **Windows**: Run `.\build.ps1`
- **Linux/macOS**: Run `make build` (Requires `make` to be installed)
