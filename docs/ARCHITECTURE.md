# System Architecture

## 1. Cryptography Pipeline

The system uses a symmetric encryption mechanism to protect the Vault. Data is never stored as plain-text.

Master Password ──> [ Argon2id + Salt + Dynamic Header Params ] ──> Master Key (32 bytes)
                                                                            │
Plaintext Data  ──> [ AES-256-GCM + Nonce ] <───────────────────────────────┘
                                                                            │
                                                                            ▼
                                                                  [ Encrypted Vault File ]

---

## 2. Storage Mechanism

- Default location: `~/.passmgr/vault.enc`
- Vault file structure format (V2 Binary with versioning support):
  - **Legacy V1 Layout**: `[Salt (16 bytes)] + [Nonce (12 bytes)] + [AES-GCM Ciphertext]`
  - **Modern V2 Layout**: `[Magic "PMV2" (4 bytes)] + [Salt (16 bytes)] + [Nonce (12 bytes)] + [Time (4 bytes)] + [Memory (4 bytes)] + [Threads (1 byte)] + [AES-GCM Ciphertext]`

---

## 3. Cobra Command Tree

```
root (passmgr)
├── init       (Initialize a new vault)
├── add        (Add a new credential)
├── get        (Retrieve and decrypt a credential)
├── list       (List all saved services)
├── delete     (Delete a credential)
├── update     (Update an existing credential)
├── search     (Search for a credential)
├── changepass (Change the master password)
├── audit      (Check password strength & data breaches)
├── import     (Import credentials from unencrypted JSON)
├── export     (Export credentials to unencrypted JSON)
├── generate   (Utility to quickly generate secure passwords)
└── tui        (Launch the interactive split-screen visual dashboard)
```
