package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

var (
	ConfigDir   = filepath.Join(os.Getenv("HOME"), ".config", "vpn-tool")
	AuthFile    = filepath.Join(ConfigDir, "auth.json")
	SessionFile = filepath.Join(ConfigDir, "session")
)

type AuthData struct {
	Version int    `json:"version"`
	PinHash string `json:"pin_hash"`
}

func ensureConfig() error {
	return os.MkdirAll(ConfigDir, 0700)
}

func hashPin(pin string) string {
	hash := sha256.Sum256([]byte(pin))
	return hex.EncodeToString(hash[:])
}

func registered() bool {
	_, err := os.Stat(AuthFile)
	return err == nil
}

func loggedIn() bool {
	data, err := os.ReadFile(SessionFile)
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(data)) == "logged_in"
}

func loadAuth() *AuthData {
	if !registered() {
		return nil
	}
	data, err := os.ReadFile(AuthFile)
	if err != nil {
		return nil
	}
	var authData AuthData
	if err := json.Unmarshal(data, &authData); err != nil {
		return nil
	}
	return &authData
}

func RequireLogin() {
	if !registered() {
		fmt.Fprintln(os.Stderr, "Not registered. Run: vpn-tool auth register")
		os.Exit(1)
	}
	if !loggedIn() {
		fmt.Fprintln(os.Stderr, "Not logged in. Run: vpn-tool auth login")
		os.Exit(1)
	}
}

func ask(prompt string) (string, error) {
	fmt.Print(prompt)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func NewAuthCmd() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands",
		Long: heredoc.Doc(`
			Manage authentication, PIN registration, and session states for vpn-tool.
			You must register a PIN before logging in and executing protected commands.
		`),
		Example: heredoc.Doc(`
			$ vpn-tool auth register
			$ vpn-tool auth login
			$ vpn-tool auth status
		`),
	}

	authCmd.AddCommand(&cobra.Command{
		Use:   "register",
		Short: "Register a new PIN",
		Long:  heredoc.Doc(`Register a new security PIN for authenticating subsequent session operations.`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if registered() {
				fmt.Fprintln(os.Stderr, "PIN already registered.")
				os.Exit(1)
			}
			pin, err := ask("Enter new PIN: ")
			if err != nil || pin == "" {
				fmt.Fprintln(os.Stderr, "PIN cannot be empty.")
				os.Exit(1)
			}
			confirm, err := ask("Confirm new PIN: ")
			if err != nil || pin != confirm {
				fmt.Fprintln(os.Stderr, "PINs do not match.")
				os.Exit(1)
			}
			if err := ensureConfig(); err != nil {
				return err
			}
			authData := AuthData{Version: 1, PinHash: hashPin(pin)}
			jsonData, err := json.MarshalIndent(authData, "", " ")
			if err != nil {
				return err
			}
			if err := os.WriteFile(AuthFile, jsonData, 0600); err != nil {
				return err
			}
			if err := os.WriteFile(SessionFile, []byte("logged_out"), 0600); err != nil {
				return err
			}
			fmt.Println("PIN registered successfully.")
			fmt.Println("Run: vpn-tool auth login")
			return nil
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "login",
		Short: "Login with PIN",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !registered() {
				fmt.Fprintln(os.Stderr, "No PIN registered. Run: vpn-tool auth register")
				os.Exit(1)
			}
			authData := loadAuth()
			pin, err := ask("Enter PIN: ")
			if err != nil {
				return err
			}
			if authData == nil || authData.PinHash != hashPin(pin) {
				fmt.Fprintln(os.Stderr, "Invalid PIN.")
				os.Exit(1)
			}
			if err := ensureConfig(); err != nil {
				return err
			}
			if err := os.WriteFile(SessionFile, []byte("logged_in"), 0600); err != nil {
				return err
			}
			fmt.Println("Login successful.")
			return nil
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "logout",
		Short: "Logout",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !registered() {
				fmt.Fprintln(os.Stderr, "No PIN registered.")
				os.Exit(1)
			}
			authData := loadAuth()
			pin, err := ask("Enter PIN: ")
			if err != nil {
				return err
			}
			if authData == nil || authData.PinHash != hashPin(pin) {
				fmt.Fprintln(os.Stderr, "Invalid PIN.")
				os.Exit(1)
			}
			if err := os.WriteFile(SessionFile, []byte("logged_out"), 0600); err != nil {
				return err
			}
			fmt.Println("Logout successful.")
			return nil
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		Run: func(cmd *cobra.Command, args []string) {
			regStatus := "no"
			if registered() {
				regStatus = "yes"
			}
			loginStatus := "no"
			if loggedIn() {
				loginStatus = "yes"
			}
			fmt.Printf("Registered: %s\n", regStatus)
			fmt.Printf("Logged in: %s\n", loginStatus)
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "reset",
		Short: "Reset PIN",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !registered() {
				fmt.Fprintln(os.Stderr, "No PIN registered.")
				os.Exit(1)
			}
			authData := loadAuth()
			current, err := ask("Enter current PIN: ")
			if err != nil {
				return err
			}
			if authData == nil || authData.PinHash != hashPin(current) {
				fmt.Fprintln(os.Stderr, "Invalid current PIN.")
				os.Exit(1)
			}
			newPin, err := ask("Enter new PIN: ")
			if err != nil || newPin == "" {
				fmt.Fprintln(os.Stderr, "New PIN cannot be empty.")
				os.Exit(1)
			}
			confirm, err := ask("Confirm new PIN: ")
			if err != nil || newPin != confirm {
				fmt.Fprintln(os.Stderr, "PINs do not match.")
				os.Exit(1)
			}
			newAuthData := AuthData{Version: 1, PinHash: hashPin(newPin)}
			jsonData, err := json.MarshalIndent(newAuthData, "", " ")
			if err != nil {
				return err
			}
			if err := os.WriteFile(AuthFile, jsonData, 0600); err != nil {
				return err
			}
			if err := os.WriteFile(SessionFile, []byte("logged_out"), 0600); err != nil {
				return err
			}
			fmt.Println("PIN reset successfully.")
			return nil
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "remove",
		Short: "Remove authentication",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !registered() {
				fmt.Fprintln(os.Stderr, "No PIN registered.")
				os.Exit(1)
			}
			authData := loadAuth()
			pin, err := ask("Enter PIN: ")
			if err != nil {
				return err
			}
			if authData == nil || authData.PinHash != hashPin(pin) {
				fmt.Fprintln(os.Stderr, "Invalid PIN.")
				os.Exit(1)
			}
			_ = os.Remove(AuthFile)
			_ = os.Remove(SessionFile)
			fmt.Println("Authentication removed.")
			return nil
		},
	})

	return authCmd
}
