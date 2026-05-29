# Features Specification

## 1. `init` Command

- **Behavior:** Prompts the user to enter a Master Password (hides characters during input).
- **Processing:** Generates a random Salt, initializes an empty vault, and saves it securely to the file system.

## 2. `add <service_name>` Command

- **Behavior:** Receives the service name from the argument. Prompts for the username and password to be saved (or auto-generates if `-g` is used).
- **Processing:** Checks if the vault file exists. If so, decrypts it, inserts the new entry, re-encrypts the entirety, and overwrites the file securely.

## 3. `get <service_name>` Command

- **Behavior:** Prompts for the Master Password, decrypts the corresponding entry, and prints it to the screen.
- **Options:** Supports the `-c` or `--copy` flag to copy directly to the clipboard instead of printing to the terminal.

## 4. `list` Command

- **Behavior:** Prints out all saved services (without their passwords) so the user can see what's in the vault.

## 5. `generate` Command

- **Behavior:** Generates a secure random string based on `crypto/rand`.
- **Options:** Customizable length. Can omit special characters using `-n`.
