package vault

import (
	"bytes"
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

	searchStr := []byte(suffix + ":")
	if bytes.Contains(bodyBytes, searchStr) {
		return true, nil
	}

	return false, nil
}

// RunAudit performs a full audit on the vault and returns an AuditReport.
func RunAudit(vault *Vault, online bool) AuditReport {
	var report AuditReport
	report.TotalChecked = len(vault.Entries)

	// Check for reuse
	pwMap := make(map[string][]string)
	for service, entry := range vault.Entries {
		pwMap[entry.Password] = append(pwMap[entry.Password], service)
	}

	for _, services := range pwMap {
		if len(services) > 1 {
			report.ReusedCount++
			report.Messages = append(report.Messages, AuditMessage{
				Service: services[0], // Use the first service as representative, or customize later
				Level:   "REUSED",
				Message: fmt.Sprintf("The password for %s is used across %d services.", services[0], len(services)),
			})
		}
	}

	for service, entry := range vault.Entries {
		// Strength check using shared crypto function
		_, isWeak := EvaluatePasswordStrength(entry.Password)
		if isWeak {
			report.WeakCount++
			report.Messages = append(report.Messages, AuditMessage{
				Service: service,
				Level:   "WEAK",
				Message: fmt.Sprintf("Password for '%s' is weak (under 8 chars or lacks complexity).", service),
			})
		}

		// Online HIBP check
		if online {
			pwned, err := CheckPwned(entry.Password)
			if err != nil {
				report.Messages = append(report.Messages, AuditMessage{
					Service: service,
					Level:   "ERROR",
					Message: fmt.Sprintf("HIBP API Check failed for '%s': %v", service, err),
				})
			} else if pwned {
				report.PwnedCount++
				report.Messages = append(report.Messages, AuditMessage{
					Service: service,
					Level:   "PWNED",
					Message: fmt.Sprintf("Password for '%s' has been found in data breaches!", service),
				})
			}
		}
	}

	return report
}
