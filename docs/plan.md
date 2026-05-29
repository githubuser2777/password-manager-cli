# Technical Plan: Password Manager CLI Enhancements

## 1. Updated Architecture & File Layout
The core layout remains aligned with the Standard Go Project Layout. The major modules affected:
- `internal/crypto/crypto.go`: Implement secure random shuffling, guaranteed char classes, memory zeroing, Argon2 parameter dynamic headers, master password strength check.
- `internal/storage/storage.go`: Handle V1 (Legacy) and V2 (Dynamic Header with Magic Bytes `"PMV2"`) formats. Zero out JSON plaintext and key slices.
- `internal/tui/`: Introduce a revamped dashboard TUI layout (split-screen, responsive, dynamic detail view), dynamic password strength indicator, TUI unlocking spinner, and background clipboard auto-clearing via `tea.Tick`.
- `cmd/`: Integrate memory zeroing, prompt adjustments, master password strength checks on initialization/change, and interactive CLI countdown for clipboard clearing.

---

## 2. Technical Designs

### A. Dynamic Vault Format & Backward Compatibility
To avoid breaking any existing vault file, we define a hybrid storage format.
- **Magic Signature**: The first 4 bytes of a V2 file are the ASCII bytes of `"PMV2"` (`[0x50, 0x4D, 0x56, 0x32]`).
- **Layout Comparisons**:
  - **V1 (Legacy)**: `[Salt (16 bytes)] + [Nonce (12 bytes)] + [AES-GCM Ciphertext]`
    - Key derivation parameters are hardcoded to standard defaults (Time=1, Memory=64MB, Threads=4).
  - **V2 (Modern)**: `[Magic "PMV2" (4 bytes)] + [Salt (16 bytes)] + [Nonce (12 bytes)] + [Time (4 bytes, BigEndian)] + [Memory (4 bytes, BigEndian)] + [Threads (1 byte)] + [AES-GCM Ciphertext]`
    - Key derivation parameters are read dynamically from the header.

### B. Memory Security & Hardening
Go uses an automatic garbage collector, and strings are immutable, making them hard to overwrite. We transition to using `[]byte` for critical components:
1. `ZeroBytes(b []byte)` function to clear sensitive data:
   ```go
   func ZeroBytes(b []byte) {
       for i := range b {
           b[i] = 0
       }
   }
   ```
2. Prompt functions returning `[]byte` instead of `string`.
3. Wipe/Zero all derived keys, master passwords, and JSON plaintext slices as soon as decryption/encryption is complete.

### C. Master Password Strength Validator
Implement in `internal/crypto`:
- **Rules**: Length >= 12, contains uppercase, lowercase, digit, and special character.

### D. Guaranteed Random Password Generation
To ensure generated passwords are valid for any site:
1. Pick one random character from Lowercase, Uppercase, Numbers, and Special (if enabled).
2. Generate the remaining characters randomly from the permitted pool.
3. Perform a cryptographically secure Fisher-Yates shuffle using `crypto/rand` to randomize position.

### E. Interactive CLI Clipboard Countdown
Instead of doing a silent, blocking `time.Sleep`, `passmgr get -c` will:
1. Write password to clipboard.
2. Render a dynamic terminal countdown bar (e.g. `[██████████████████] 30s`).
3. Refresh every 1 second (or 100ms for smooth bar movement).
4. Catch interrupt signals (Ctrl+C) to gracefully clear the clipboard immediately and exit.

### F. Premium Revamped TUI Dashboard & Dynamic UX
1. **Split-Screen Dashboard Layout (`stateList`)**:
   - Split the screen horizontally into a left and right panel using `lipgloss.JoinHorizontal`:
     - **Left Panel (List)**: Shows the interactive list of services.
     - **Right Panel (Details Card)**: Beautifully renders the selected item's username, masked password ("••••••••"), notes, dates, and action keys. A double-border styling with custom padding will wrap this details panel.
     - **Responsive Layout**: If the window width is smaller than 80 characters, automatically fall back to showing only the list, switching to a full-screen view mode when an item is selected.
2. **Color Palette Tokens**:
   - Accent Purple/Lavender: `#BD93F9` / `Color("141")`
   - Accent Highlight (Orchid/Pink): `#FF79C6` / `Color("205")`
   - Dark Background Slate: `#282A36` / `Color("236")`
   - Mint Green (Strong/Success): `#50FA7B` / `Color("120")`
   - Amber (Medium/Warning): `#F1FA8C` / `Color("228")`
   - Coral (Weak/Danger): `#FF5555` / `Color("203")`
3. **Dynamic Password Strength Meter**:
   - In the Add/Edit form, as the user types into the Password input box, the TUI dynamically calculates password strength.
   - Render a color-coded bar below the field:
     - Weak: `Strength: [■■□□□□] Weak` (Coral Red)
     - Medium: `Strength: [■■■■□□] Medium` (Amber)
     - Strong: `Strength: [■■■■■■] Strong` (Mint Green)
4. **Unlocking Spinner (`stateDecrypting`)**:
   - When the user submits the Master Password, show a loader screen using `github.com/charmbracelet/bubbles/spinner` with custom color transitions.
   - Spin off a Bubble Tea command performing key derivation asynchronously, preventing thread hang.
5. **Non-Blocking Clipboard Timer Tick**:
   - In TUI, copying credentials triggers a background tick counting down from 30 seconds. Display remaining time at the bottom status bar and zero out clipboard safely upon expiration.
