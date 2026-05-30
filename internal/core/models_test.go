package core

import (
	"encoding/json"
	"testing"
)

func TestVaultJSONAndStructIntegrity(t *testing.T) {
	v := Vault{
		Salt: []byte("test_salt"),
		Entries: map[string]Entry{
			"test_site": {
				Username:  "user1",
				Password:  "pass1",
				Notes:     "some notes",
				CreatedAt: "2023-01-01",
				UpdatedAt: "2023-01-02",
			},
		},
	}

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to marshal vault: %v", err)
	}

	var v2 Vault
	if err := json.Unmarshal(data, &v2); err != nil {
		t.Fatalf("Failed to unmarshal vault: %v", err)
	}

	if string(v2.Salt) != "test_salt" {
		t.Errorf("Expected salt 'test_salt', got '%s'", string(v2.Salt))
	}

	entry, ok := v2.Entries["test_site"]
	if !ok {
		t.Fatalf("Expected entry 'test_site' to exist")
	}

	if entry.Username != "user1" || entry.Password != "pass1" || entry.Notes != "some notes" {
		t.Errorf("Entry details mismatch: %+v", entry)
	}
}
