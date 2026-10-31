package test

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

func NewTestCmd() *cobra.Command {
	var ms string

	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test API delay",
		Long: heredoc.Doc(`
			Test the API latency/delay by requesting a sleep duration on the server.
			Requires an active authenticated session.
		`),
		Example: heredoc.Doc(`
			$ vpn-tool test --ms 200
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			auth.RequireLogin()
			endpoint := fmt.Sprintf("/api/vpn/test/sleep?ms=%s", url.QueryEscape(ms))
			return apiRequest("GET", endpoint)
		},
	}

	cmd.Flags().StringVar(&ms, "ms", "100", "Milliseconds")

	return cmd
}
