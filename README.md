# passmgr

[![License: GPL-3.0-only](https://img.shields.io/badge/License-GPL--3.0--only-blue.svg)](LICENSE)

A fast, secure, and minimal command-line password manager featuring a local vault and an interactive Terminal UI.

Secured by industry-standard cryptography (AES-256-GCM and Argon2id).

## Features

- **Local-First**: Your vault is stored locally at `~/.passmgr/vault.enc`. No cloud dependency.
- **Strong Crypto**: AES-256-GCM authenticated encryption with Argon2id for key derivation.
- **Interactive TUI**: Fast, keyboard-driven interface for daily use.
- **Security Audit**: Built-in checks for weak, reused, or compromised passwords (via HaveIBeenPwned).
- **Zero Lock-in**: Easy JSON import and export capabilities.

## Installation

Requires Go 1.21 or later.

```sh
git clone https://github.com/githubuser2777/password-manager-cli.git
cd password-manager-cli

# Global install (accessible anywhere as password-manager-cli)
go install

# Or build locally (creates passmgr in current directory)
go build -o passmgr      # Linux/macOS
go build -o passmgr.exe  # Windows
```

## Quick Start

*Note: If you used `go install` above, your command will be `password-manager-cli` everywhere.*

1. Initialize your vault:
   ```sh
   .\passmgr.exe init  # Windows (Local build)
   ./passmgr init      # Linux/macOS (Local build)
   ```
   > **Note:** Creates `~/.passmgr/vault.enc` and sets your master password. Strength is enforced. No cloud recovery exists—if you lose it, your vault is gone.

2. Launch the interactive TUI:
   ```sh
   .\passmgr.exe   # Windows
   ./passmgr       # Linux/macOS
   ```

## Usage

### Interactive TUI

Simply run `passmgr` without arguments to enter the interactive mode.

| Key | Action |
| --- | --- |
| `a` | Add entry |
| `e` | Edit entry |
| `d` | Delete entry |
| `c` | Copy password to clipboard |
| `Enter` | View details |
| `/` | Search |
| `r` | Run local security audit |
| `R` | Run HaveIBeenPwned audit |

### CLI Commands

All features are also available directly from the command line:

```sh
passmgr add <domain> [-g]    # Add entry (-g to auto-generate password)
passmgr get <domain> [-c]    # View entry (-c to copy password)
passmgr update <domain>      # Edit entry
passmgr delete <domain>      # Delete entry
passmgr search <query>       # Search entries
passmgr list                 # List all domains
passmgr generate <len>       # Generate a random password
passmgr audit                # Run security audit
passmgr export               # Export vault to JSON
passmgr import               # Import vault from JSON
passmgr changepass           # Change master password
```

## License

This project is licensed under the GNU General Public License v3.0 (GPL-3.0-only). See [`LICENSE`](LICENSE) for details.

<!-- ponytail: Professional doesn't mean bloated. Added explicit constraints to 'init' so users know what to expect (no recovery). No unnecessary fluff. YAGNI. -->
