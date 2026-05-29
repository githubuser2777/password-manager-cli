package core

// Entry represents a single credential stored in the vault.
type Entry struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Notes     string `json:"notes,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// Vault represents the entire collection of credentials.
type Vault struct {
	Salt    []byte           `json:"salt"`
	Entries map[string]Entry `json:"entries"`
}
