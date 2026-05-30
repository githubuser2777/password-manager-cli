package core

import (
	"testing"
)

func TestCheckPwned(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping network test in short mode")
	}

	// 'password' is a known pwned password
	isPwned, err := CheckPwned("password")
	if err != nil {
		t.Fatalf("CheckPwned failed: %v", err)
	}
	if !isPwned {
		t.Errorf("Expected 'password' to be marked as pwned")
	}

	// A highly randomized string should not be pwned
	isPwned, err = CheckPwned("kjshdkfjhskdjhfksdjhfksjdfh123981293812")
	if err != nil {
		t.Fatalf("CheckPwned failed: %v", err)
	}
	if isPwned {
		t.Errorf("Expected highly randomized string to not be pwned")
	}
}
