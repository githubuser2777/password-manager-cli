# Source Code Mapping (Context)

## 1. Module Locations

- `internal/crypto/`: Contains Argon2id hashing and AES-GCM encryption logic. If the encryption algorithm needs modification, only intervene in this area.
- `internal/storage/`: Contains reading/writing logic interacting with the file system securely (Atomic Write).
- `internal/core/`: Contains the domain models (`Vault`, `Entry`).
- `cmd/`: Contains definitions for the Cobra CLI interface.

## 2. Read/Write Rules for Agents

- When a feature change is requested: Update specifications first, then proceed to modify the code in the `cmd/` directory.
- When optimizing performance or security: Cross-check with `.github/SECURITY.md` to ensure core principles are not violated.
