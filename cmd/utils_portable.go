//go:build portable

package cmd

import (
	"os"
	"path/filepath"
)

func init() {
	// Override getVaultPath to store the vault in the executable's directory
	getVaultPath = func() string {
		exe, err := os.Executable()
		if err != nil {
			return "vault.enc"
		}
		return filepath.Join(filepath.Dir(exe), "vault.enc")
	}
}
