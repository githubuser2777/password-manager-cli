# Security Policy

## Supported Versions

Currently, the `master` branch is the only supported version. If you are using this CLI, please make sure to pull the latest changes.

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability within Password Manager CLI, please do NOT open a public issue. Instead, please send an email to the repository owner or use GitHub's private vulnerability reporting feature.

### Security Mechanisms in this Project:
- **AES-256-GCM**: Used for authenticated encryption of the vault.
- **Argon2id**: Used for key derivation to prevent brute-force attacks on the Master Password.
- **Memory Management**: The master password is not persistently stored, and the vault file requires `chmod 600` local permissions.

Thank you for helping keep this project secure!
