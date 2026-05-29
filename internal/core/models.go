package core

// Entry represents a single credential stored in the vault.
type Entry struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

// Vault represents the entire collection of credentials.
type Vault struct {
	Salt    []byte           `json:"salt"`
	Entries map[string]Entry `json:"entries"`
}
