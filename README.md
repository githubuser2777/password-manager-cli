# passmgr

A CLI password manager with a local vault. Secured by AES-256-GCM and Argon2id.

## Installation

Requires Go 1.21 or later.

```sh
go build -o passmgr
```

## Quick Start

```sh
passmgr init  # Initialize vault at ~/.passmgr/vault.enc
passmgr       # Open the interactive TUI
```

## Usage

### Interactive TUI

Run `passmgr` without arguments to enter the TUI.

* `a`: Add entry
* `e`: Edit entry
* `d`: Delete entry
* `c`: Copy password
* `Enter`: View details
* `/`: Search
* `r` / `R`: Security audit (Local / HaveIBeenPwned)

### CLI Commands

```sh
passmgr init
passmgr add <domain> [-g]    # Add entry (-g to auto-generate)
passmgr get <domain> [-c]    # View entry (-c to copy password)
passmgr update <domain>      # Edit entry
passmgr delete <domain>      # Delete entry
passmgr search <query>       # Search entries
passmgr list                 # List all domains
passmgr generate <len>       # Generate random password
passmgr audit                # Run security audit
passmgr export               # Export vault (JSON)
passmgr import               # Import vault (JSON)
passmgr changepass           # Change master password
```

<!-- ponytail: Kept README strictly technical. Added clearer structure (Installation, TUI vs CLI layout). No badges, no fluff, no unnecessary sections. YAGNI. -->
