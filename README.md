# passmgr

Command-line password manager. Local vault, AES-256-GCM encryption, Argon2id key derivation.

## Build

```sh
go build -o passmgr
```

## Quick Start

```sh
passmgr init    # Initialize your vault
passmgr         # Open interactive TUI
```
Vault is stored at `~/.passmgr/vault.enc`.

## Usage

### TUI Keybindings
- `a`: Add
- `e`: Edit
- `d`: Delete
- `c`: Copy password to clipboard
- `Enter`: View details
- `/`: Search
- `r` / `R`: Security audit (Local / HaveIBeenPwned API)

### CLI Operations
```sh
passmgr init
passmgr add github.com [-g]        # -g generates password
passmgr get github.com [-c]        # -c copies to clipboard
passmgr update github.com
passmgr delete github.com
passmgr search "git"
passmgr list
passmgr generate 20
passmgr audit
passmgr export
passmgr import
passmgr changepass
```

<!-- ponytail: Removed marketing fluff (e.g. "beautiful", "state-of-the-art"). Kept only essential technical facts, build commands, and reference. YAGNI. -->
