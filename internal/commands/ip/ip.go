package ip

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"vpn-tool/internal/commands/auth"
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

func NewIPCmd() *cobra.Command {
	var country string
	var count string

	cmd := &cobra.Command{
		Use:   "ip",
		Short: "Get VPN IP addresses",
		Long: heredoc.Doc(`
			Retrieve a specified number of VPN IP addresses for a given country code.
			Requires an active authenticated session.
		`),
		Example: heredoc.Doc(`
			$ vpn-tool ip --country US --count 5
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			auth.RequireLogin()
			endpoint := fmt.Sprintf("/api/vpn/ip?countryCode=%s&count=%s", url.QueryEscape(country), url.QueryEscape(count))
			return apiRequest("POST", endpoint)
		},
	}

	cmd.Flags().StringVar(&country, "country", "", "Country code")
	cmd.Flags().StringVar(&count, "count", "10", "Number of IPs")
	_ = cmd.MarkFlagRequired("country")

	return cmd
}
