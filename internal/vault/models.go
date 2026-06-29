package vault

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

// AuditMessage represents a single warning from the audit.
type AuditMessage struct {
	Service string
	Level   string // "WEAK", "REUSED", "PWNED", "ERROR"
	Message string
}

// AuditReport holds the full results of an audit.
type AuditReport struct {
	TotalChecked int
	WeakCount    int
	ReusedCount  int
	PwnedCount   int
	Messages     []AuditMessage
}
