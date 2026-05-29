# System Architecture

## 1. Cryptography Pipeline

The system uses a symmetric encryption mechanism to protect the Vault. Data is never stored as plain-text.

Master Password ──> [ Argon2id + Salt ] ──> Master Key (32 bytes)
                                                  │
Plaintext Data  ──> [ AES-256-GCM + Nonce ] <──────┘
                                                  │
                                                  ▼
                                        [ Encrypted Vault File ]

## 2. Storage Mechanism

- Default location: `~/.passmgr/vault.enc`
- Vault file structure format (Binary):
  `Salt (16 bytes) + Nonce (12 bytes) + Ciphertext`

## 3. Cobra Command Tree

root (passmgr)
├── init (Initialize the vault)
├── add  (Add a new account)
├── get  (Retrieve and decrypt)
├── list (List all saved services)
└── generate (Generate a random password)
