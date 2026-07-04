package sys

import (
	"bytes"
	"os/exec"
	"runtime"
)

// WriteClipboard writes text to the system clipboard using native commands.
// ponytail: Standard library os/exec replaces github.com/atotto/clipboard dependency.
func WriteClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("clip")
	case "darwin":
		cmd = exec.Command("pbcopy")
	default: // linux / bsd
		// Try wl-copy (Wayland) first, fallback to xclip (X11)
		if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd = exec.Command("wl-copy")
		} else {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		}
	}

	cmd.Stdin = bytes.NewBufferString(text)
	return cmd.Run()
}
