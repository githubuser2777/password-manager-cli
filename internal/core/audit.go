package core

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CheckPwned securely checks the HaveIBeenPwned API to see if a password has been leaked.
func CheckPwned(password string) (bool, error) {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	hashStr := strings.ToUpper(fmt.Sprintf("%x", hasher.Sum(nil)))

	prefix := hashStr[:5]
	suffix := hashStr[5:]

	resp, err := http.Get("https://api.pwnedpasswords.com/range/" + prefix)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(bodyBytes), "\n")
	for _, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), ":")
		if len(parts) >= 1 && parts[0] == suffix {
			return true, nil // Found match
		}
	}

	return false, nil
}
