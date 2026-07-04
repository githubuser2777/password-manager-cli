package cmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"password-manager-cli/internal/sys"
	"password-manager-cli/internal/tui"
	"password-manager-cli/internal/vault"

	"golang.org/x/term"
)

var getVaultPath = func() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting home directory:", err)
		os.Exit(1)
	}
	return filepath.Join(home, ".passmgr", "vault.enc")
}

var promptPassword = func(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return nil, err
	}
	trimmed := []byte(strings.TrimSpace(string(bytePassword)))
	vault.ZeroBytes(bytePassword)
	return trimmed, nil
}

func withVault(action func(v *vault.Vault, masterPw []byte, path string)) {
	path := getVaultPath()
	masterPw, err := promptPassword("Master Password: ")
	if err != nil {
		return
	}
	defer vault.ZeroBytes(masterPw)

	v, err := vault.LoadVault(path, masterPw)
	if err != nil {
		fmt.Println("Failed to open vault:", err)
		return
	}
	action(v, masterPw, path)
}

func withVaultMutate(action func(v *vault.Vault, masterPw []byte, path string) bool) {
	withVault(func(v *vault.Vault, masterPw []byte, path string) {
		if action(v, masterPw, path) {
			if err := vault.SaveVault(path, masterPw, v); err != nil {
				fmt.Println("Failed to save vault:", err)
			}
		}
	})
}

// Execute parses arguments and routes to the appropriate command.
// ponytail: Replaced the entire Cobra abstraction and 16 files with standard library flags.
func Execute() {
	if len(os.Args) < 2 {
		if err := tui.StartApp(getVaultPath()); err != nil {
			fmt.Println("Error running TUI:", err)
			os.Exit(1)
		}
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "tui":
		if err := tui.StartApp(getVaultPath()); err != nil {
			fmt.Println("Error running TUI:", err)
			os.Exit(1)
		}
	case "init":
		path := getVaultPath()
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Vault already exists at", path)
			return
		}
		masterPw, err := promptPassword("Set Master Password: ")
		if err != nil {
			return
		}
		defer vault.ZeroBytes(masterPw)

		confirmPw, err := promptPassword("Confirm Master Password: ")
		if err != nil {
			return
		}
		defer vault.ZeroBytes(confirmPw)

		if string(masterPw) != string(confirmPw) {
			fmt.Println("Passwords do not match.")
			return
		}

		if err := vault.ValidateMasterPassword(masterPw); err != nil {
			fmt.Println("Weak password:", err)
			return
		}

		salt, err := vault.GenerateSalt(16)
		if err != nil {
			fmt.Println("Failed to generate salt:", err)
			return
		}

		v := &vault.Vault{
			Salt:    salt,
			Entries: make(map[string]vault.Entry),
		}

		if err := vault.SaveVault(path, masterPw, v); err != nil {
			fmt.Println("Failed to initialize vault:", err)
		} else {
			fmt.Println("Vault initialized successfully at", path)
		}

	case "add":
		fs := flag.NewFlagSet("add", flag.ExitOnError)
		generateFlag := fs.Bool("g", false, "Auto-generate password")
		fs.Parse(args)

		if fs.NArg() != 1 {
			fmt.Println("Usage: passmgr add <service> [-g]")
			return
		}
		service := fs.Arg(0)

		withVaultMutate(func(v *vault.Vault, masterPw []byte, path string) bool {
			if _, exists := v.Entries[service]; exists {
				fmt.Printf("Service '%s' already exists in the vault.\n", service)
				return false
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Username: ")
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)

			var password string
			if *generateFlag {
				pw, err := vault.GenerateRandomPassword(16, true)
				if err != nil {
					fmt.Println("Error generating password:", err)
					return false
				}
				password = pw
				fmt.Println("Generated Password:", password)
			} else {
				rawPw, err := promptPassword("Password: ")
				if err != nil {
					return false
				}
				password = string(rawPw)
				vault.ZeroBytes(rawPw)
			}

			fmt.Print("Notes (optional): ")
			notes, _ := reader.ReadString('\n')
			notes = strings.TrimSpace(notes)

			now := time.Now().Format(time.RFC3339)
			v.Entries[service] = vault.Entry{
				Username:  username,
				Password:  password,
				Notes:     notes,
				CreatedAt: now,
				UpdatedAt: now,
			}

			fmt.Printf("Successfully added '%s' to the vault.\n", service)
			return true
		})

	case "get":
		fs := flag.NewFlagSet("get", flag.ExitOnError)
		copyFlag := fs.Bool("c", false, "Copy password to clipboard")
		fs.Parse(args)

		if fs.NArg() != 1 {
			fmt.Println("Usage: passmgr get <service> [-c]")
			return
		}
		service := fs.Arg(0)

		withVault(func(v *vault.Vault, masterPw []byte, path string) {
			entry, exists := v.Entries[service]
			if !exists {
				fmt.Printf("Service '%s' not found.\n", service)
				return
			}
			fmt.Printf("Service: %s\nUsername: %s\nPassword: %s\nNotes: %s\n", service, entry.Username, entry.Password, entry.Notes)

			if *copyFlag {
				if err := sys.WriteClipboard(entry.Password); err != nil {
					fmt.Println("Failed to copy to clipboard:", err)
				} else {
					fmt.Println("Password copied to clipboard.")
				}
			}
		})

	case "update":
		if len(args) != 1 {
			fmt.Println("Usage: passmgr update <service>")
			return
		}
		service := args[0]
		withVaultMutate(func(v *vault.Vault, masterPw []byte, path string) bool {
			entry, exists := v.Entries[service]
			if !exists {
				fmt.Printf("Service '%s' not found.\n", service)
				return false
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("New Username (leave blank to keep '%s'): ", entry.Username)
			newUsername, _ := reader.ReadString('\n')
			newUsername = strings.TrimSpace(newUsername)
			if newUsername != "" {
				entry.Username = newUsername
			}

			rawPw, err := promptPassword("New Password (leave blank to keep current): ")
			if err == nil && len(rawPw) > 0 {
				entry.Password = string(rawPw)
				vault.ZeroBytes(rawPw)
			}

			fmt.Printf("New Notes (leave blank to keep '%s'): ", entry.Notes)
			newNotes, _ := reader.ReadString('\n')
			newNotes = strings.TrimSpace(newNotes)
			if newNotes != "" {
				entry.Notes = newNotes
			}

			entry.UpdatedAt = time.Now().Format(time.RFC3339)
			v.Entries[service] = entry
			fmt.Printf("Successfully updated '%s'.\n", service)
			return true
		})

	case "delete":
		if len(args) != 1 {
			fmt.Println("Usage: passmgr delete <service>")
			return
		}
		service := args[0]
		withVaultMutate(func(v *vault.Vault, masterPw []byte, path string) bool {
			if _, exists := v.Entries[service]; !exists {
				fmt.Printf("Service '%s' not found.\n", service)
				return false
			}

			fmt.Printf("Are you sure you want to delete '%s'? (y/N): ", service)
			reader := bufio.NewReader(os.Stdin)
			confirm, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
				fmt.Println("Aborted.")
				return false
			}

			delete(v.Entries, service)
			fmt.Printf("Successfully deleted '%s'.\n", service)
			return true
		})

	case "list":
		withVault(func(v *vault.Vault, masterPw []byte, path string) {
			if len(v.Entries) == 0 {
				fmt.Println("Vault is empty.")
				return
			}
			fmt.Println("Saved Services:")
			for service := range v.Entries {
				fmt.Println("-", service)
			}
		})

	case "search":
		if len(args) != 1 {
			fmt.Println("Usage: passmgr search <query>")
			return
		}
		query := strings.ToLower(args[0])
		withVault(func(v *vault.Vault, masterPw []byte, path string) {
			fmt.Println("Search Results:")
			found := false
			for service := range v.Entries {
				if strings.Contains(strings.ToLower(service), query) {
					fmt.Println("-", service)
					found = true
				}
			}
			if !found {
				fmt.Println("No matches found.")
			}
		})

	case "generate":
		length := 16
		if len(args) > 0 {
			if l, err := strconv.Atoi(args[0]); err == nil && l > 0 {
				length = l
			}
		}
		pw, err := vault.GenerateRandomPassword(length, true)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(pw)

	case "audit":
		withVault(func(v *vault.Vault, masterPw []byte, path string) {
			fmt.Println("Running local security audit...")
			report := vault.RunAudit(v, false)
			for _, msg := range report.Messages {
				fmt.Printf("[%s] %s: %s\n", msg.Level, msg.Service, msg.Message)
			}
			fmt.Printf("\nAudit Complete: %d Checked | %d Weak | %d Reused\n", report.TotalChecked, report.WeakCount, report.ReusedCount)
		})

	case "export":
		withVault(func(v *vault.Vault, masterPw []byte, path string) {
			data, err := json.MarshalIndent(v.Entries, "", "  ")
			if err != nil {
				fmt.Println("Export failed:", err)
				return
			}
			fmt.Println(string(data))
		})

	case "import":
		if len(args) != 1 {
			fmt.Println("Usage: passmgr import <file.json>")
			return
		}
		file := args[0]
		withVaultMutate(func(v *vault.Vault, masterPw []byte, path string) bool {
			data, err := os.ReadFile(file)
			if err != nil {
				fmt.Println("Failed to read file:", err)
				return false
			}
			var entries map[string]vault.Entry
			if err := json.Unmarshal(data, &entries); err != nil {
				fmt.Println("Failed to parse JSON:", err)
				return false
			}
			for k, val := range entries {
				v.Entries[k] = val
			}
			fmt.Printf("Imported %d entries.\n", len(entries))
			return true
		})

	case "changepass":
		withVaultMutate(func(v *vault.Vault, masterPw []byte, path string) bool {
			newPw, err := promptPassword("New Master Password: ")
			if err != nil {
				return false
			}
			defer vault.ZeroBytes(newPw)

			confirmPw, err := promptPassword("Confirm New Master Password: ")
			if err != nil {
				return false
			}
			defer vault.ZeroBytes(confirmPw)

			if string(newPw) != string(confirmPw) {
				fmt.Println("Passwords do not match.")
				return false
			}

			if err := vault.ValidateMasterPassword(newPw); err != nil {
				fmt.Println("Weak password:", err)
				return false
			}

			newSalt, err := vault.GenerateSalt(16)
			if err != nil {
				fmt.Println("Failed to generate salt:", err)
				return false
			}

			v.Salt = newSalt
			if err := vault.SaveVault(path, newPw, v); err != nil {
				fmt.Println("Failed to save with new password:", err)
			} else {
				fmt.Println("Master password changed successfully.")
			}
			return false
		})

	case "help", "--help", "-h":
		fmt.Println("Usage: passmgr <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  init          Initialize vault")
		fmt.Println("  add           Add entry (-g to auto-generate password)")
		fmt.Println("  get           View entry (-c to copy password)")
		fmt.Println("  update        Edit entry")
		fmt.Println("  delete        Delete entry")
		fmt.Println("  list          List all domains")
		fmt.Println("  search        Search entries")
		fmt.Println("  generate      Generate random password")
		fmt.Println("  audit         Run security audit")
		fmt.Println("  export        Export vault to JSON")
		fmt.Println("  import        Import vault from JSON")
		fmt.Println("  changepass    Change master password")
		fmt.Println("  tui           Launch terminal user interface (default)")

	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Run 'passmgr help' for usage.")
		os.Exit(1)
	}
}
