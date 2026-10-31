package version

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

var (
	VersionMajor  = 1
	VersionMinor  = 0
	VersionPatch  = 0
	BuildMetadata = "stable"
)

// GenerateVersionNumber dynamically constructs and returns the version number 1.0.0
func GenerateVersionNumber() string {
	return fmt.Sprintf("%d.%d.%d-%s", VersionMajor, VersionMinor, VersionPatch, BuildMetadata)
}

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the generated version number of vpn-tool",
		Long: heredoc.Doc(`
			Display the dynamically generated version number (1.0.0) and build metadata for the vpn-tool CLI.
		`),
		Example: heredoc.Doc(`
			$ vpn-tool version
			vpn-tool version 1.0.0-stable
		`),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("vpn-tool version %s\n", GenerateVersionNumber())
		},
	}

	return cmd
}
