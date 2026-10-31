package commands

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/im7548631-design/vpn-tool/internal/commands/api"
	"github.com/im7548631-design/vpn-tool/internal/commands/auth"
	"github.com/im7548631-design/vpn-tool/internal/commands/completion"
	"github.com/im7548631-design/vpn-tool/internal/commands/connect"
	"github.com/im7548631-design/vpn-tool/internal/commands/countries"
	"github.com/im7548631-design/vpn-tool/internal/commands/create"
	"github.com/im7548631-design/vpn-tool/internal/commands/disconnect"
	"github.com/im7548631-design/vpn-tool/internal/commands/healthz"
	"github.com/im7548631-design/vpn-tool/internal/commands/ip"
	"github.com/im7548631-design/vpn-tool/internal/commands/status"
	"github.com/im7548631-design/vpn-tool/internal/commands/test"
	"github.com/im7548631-design/vpn-tool/internal/commands/version"
)

// NewRootCmd creates and returns the root cobra command for vpn-tool
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "vpn-tool",
		Short: "vpn-tool is a CLI for managing VPN gateways and secure connections",
		Long: heredoc.Doc(`
			A powerful command-line interface tool for managing authentication,
			VPN gateway instances, countries, and direct API communications.
		`),
		Example: heredoc.Doc(`
			$ vpn-tool auth register
			$ vpn-tool auth login
			$ vpn-tool countries
			$ vpn-tool create --country US
			$ vpn-tool connect --country US
		`),
	}

	// Enable automatic --version flag support using the generated version number
	rootCmd.Version = version.GenerateVersionNumber()

	// Register all subcommands
	rootCmd.AddCommand(auth.NewAuthCmd())
	rootCmd.AddCommand(healthz.NewHealthzCmd())
	rootCmd.AddCommand(api.NewApiCmd())
	rootCmd.AddCommand(countries.NewCountriesCmd())
	rootCmd.AddCommand(create.NewCreateCmd())
	rootCmd.AddCommand(connect.NewConnectCmd())
	rootCmd.AddCommand(disconnect.NewDisconnectCmd())
	rootCmd.AddCommand(status.NewStatusCmd())
	rootCmd.AddCommand(test.NewTestCmd())
	rootCmd.AddCommand(ip.NewIPCmd())
	rootCmd.AddCommand(completion.NewCompletionCmd())
	rootCmd.AddCommand(version.NewVersionCmd())

	return rootCmd
}

// Execute runs the root command and handles errors
func Execute() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
