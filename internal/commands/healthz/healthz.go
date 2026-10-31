package healthz

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/im7548631-design/vpn-tool/internal/commands/auth"
)

const APIBase = "https://vpn-gateway-manager--jjajajjajahsh.replit.app"

func apiRequest(method, endpoint string) error {
	url := fmt.Sprintf("%s%s", APIBase, endpoint)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "API request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		fmt.Fprintf(os.Stderr, "Error (%d): %s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}

	fmt.Println(string(body))
	return nil
}

func NewHealthzCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "healthz",
		Short: "Check API health",
		Long: heredoc.Doc(`
			Check the health status of the remote VPN gateway API backend.
			Requires an active authenticated session.
		`),
		Example: heredoc.Doc(`
			$ vpn-tool healthz
			API is healthy.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			auth.RequireLogin()
			return apiRequest("GET", "/api/healthz")
		},
	}

	return cmd
}
