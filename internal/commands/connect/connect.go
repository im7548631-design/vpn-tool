package connect

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/im7548631-design/vpn-tool/internal/commands/auth"
)

const APIBase = "https://vpn-gateway-manager--jjajajjajahsh.replit.app"

func apiRequest(method, endpoint string) error {
	fullURL := fmt.Sprintf("%s%s", APIBase, endpoint)

	req, err := http.NewRequest(method, fullURL, nil)
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

func NewConnectCmd() *cobra.Command {
	var country string

	cmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect VPN",
		Long: heredoc.Doc(`
			Connect to a VPN instance for the specified country code.
			Requires an active authenticated session.
		`),
		Example: heredoc.Doc(`
			$ vpn-tool connect --country US
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			auth.RequireLogin()
			endpoint := fmt.Sprintf("/api/vpn/connect?countryCode=%s", url.QueryEscape(country))
			return apiRequest("POST", endpoint)
		},
	}

	cmd.Flags().StringVar(&country, "country", "", "Country code")
	_ = cmd.MarkFlagRequired("country")

	return cmd
}
