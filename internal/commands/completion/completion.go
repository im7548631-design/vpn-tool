package completion

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

func NewCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script for shells",
		Long: heredoc.Doc(`
			Generate the autocompletion script for vpn-tool for the specified shell.
			To load completions in your current shell session:

			    $ source <(vpn-tool completion bash)

			To load completions for every new session, execute once:

			    # Linux (Bash)
			    $ vpn-tool completion bash > /etc/bash_completion.d/vpn-tool

			    # macOS (Zsh)
			    $ vpn-tool completion zsh > "${fpath[1]}/_vpn-tool"
		`),
		Example: heredoc.Doc(`
			$ vpn-tool completion bash
			$ vpn-tool completion zsh
			$ vpn-tool completion fish
			$ vpn-tool completion powershell
		`),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletion(os.Stdout)
			default:
				return fmt.Errorf("unsupported shell type %q; must be one of bash, zsh, fish, powershell", args[0])
			}
		},
	}

	return cmd
}
