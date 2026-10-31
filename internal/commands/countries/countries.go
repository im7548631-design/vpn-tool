package countries

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

func NewCountriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "countries",
		Short: "List VPN countries",
		Long: heredoc.Doc(`
			Retrieve and display the list of available VPN gateway countries from the remote API.
			Requires an active authenticated session.
		`),
		Example: heredoc.Doc(`
			$ vpn-tool countries
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			auth.RequireLogin()
			return apiRequest("GET", "/api/vpn/countries")
		},
	}

	return cmd
}
