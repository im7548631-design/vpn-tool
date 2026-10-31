package api

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"vpn-tool/internal/commands/auth"
)

type apiOptions struct {
	hostname            string
	requestMethod       string
	url                 string
	fields              []string
	formFields          []string
	requestHeaders      []string
	showResponseHeaders bool
	paginate            bool
	requestInputFile    string
	silent              bool
	outputFormat        string
	placeholderFlags    []string
}

// fieldFlagList and fieldFlagValue stub implementations for custom flag types
type fieldFlagList struct {
	fields *[]string
}

type fieldFlagValue struct {
	list *fieldFlagList
	raw  bool
}

func (f *fieldFlagValue) String() string {
	return ""
}

func (f *fieldFlagValue) Set(value string) error {
	*f.list.fields = append(*f.list.fields, value)
	return nil
}

func (f *fieldFlagValue) Type() string {
	if f.raw {
		return "string"
	}
	return "value"
}

// cmdutils mock or enum helper if needed
var outputEnum = []string{"json", "ndjson"}

func NewApiCmd() *cobra.Command {
	var opts apiOptions

	cmd := &cobra.Command{
		Use:   "api [<path>|<url>]",
		Short: "Make an authenticated HTTP request to the VPN gateway API",
		Long: heredoc.Doc(`
			Make an authenticated HTTP request to the API backend.
			If a path is provided instead of a full URL, it will be resolved against the API base URL.
		`),
		Example: heredoc.Doc(`
			$ vpn-tool api /api/healthz
			$ vpn-tool api --method POST /api/config -F status=active
			$ vpn-tool api --url https://vpn-gateway-manager--jjajajjajahsh.replit.app/api/healthz
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			auth.RequireLogin()

			targetURL := opts.url
			if targetURL == "" && len(args) > 0 {
				targetURL = args[0]
			}

			if targetURL == "" {
				return fmt.Errorf("accepts 1 arg(s), received 0; specify a path or use --url")
			}

			// If it's a relative path, prefix with API base
			baseURL := "https://vpn-gateway-manager--jjajajjajahsh.replit.app"
			if len(targetURL) > 0 && targetURL[0] == '/' {
				targetURL = baseURL + targetURL
			}

			req, err := http.NewRequest(opts.requestMethod, targetURL, nil)
			if err != nil {
				return err
			}

			// Apply headers
			for _, header := range opts.requestHeaders {
				// Simple header parsing logic can be added here
				req.Header.Set("Authorization", "Bearer token-stub")
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "API request failed: %v\n", err)
				os.Exit(1)
			}
			defer resp.Body.Close()

			if opts.showResponseHeaders {
				for k, vals := range resp.Header {
					for _, v := range vals {
						fmt.Printf("%s: %s\n", k, v)
					}
				}
				fmt.Println()
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !opts.silent {
				fmt.Println(string(body))
			}

			return nil
		},
	}

	fl := cmd.Flags()

	fl.StringVar(&opts.hostname, "hostname", "", "The hostname for the request. Defaults to gitlab.com, or the authenticated host in the current directory.")
	fl.StringVar(&opts.url, "url", "", "The full URL for the request.")
	fl.StringVarP(&opts.requestMethod, "method", "X", "GET", "The HTTP method for the request.")
	
	fieldList := &fieldFlagList{fields: &opts.fields}
	fl.VarP(&fieldFlagValue{list: fieldList}, "field", "F", "Add a parameter of inferred type. Using this flag changes the default HTTP method to POST.")
	fl.VarP(&fieldFlagValue{list: fieldList, raw: true}, "raw-field", "f", "Add a string parameter.")
	
	fl.StringArrayVar(&opts.formFields, "form", nil, "Add a multipart form field. To upload a file, prefix the value with @ followed by the file path. To read from standard input, use @- (at most once). Using this flag changes the default HTTP method to POST.")
	fl.StringArrayVarP(&opts.requestHeaders, "header", "H", nil, "Add an additional HTTP request header.")
	fl.BoolVarP(&opts.showResponseHeaders, "include", "i", false, "Include HTTP response headers in the output.")
	fl.BoolVar(&opts.paginate, "paginate", false, "Make additional HTTP requests to fetch all pages of results.")
	fl.StringVar(&opts.requestInputFile, "input", "", "The file to use as the body for the HTTP request. Sent as 'application/json' unless you set another 'Content-Type' with --header.")
	fl.BoolVar(&opts.silent, "silent", false, "Do not print the response body.")
	fl.StringVar(&opts.outputFormat, "output", "json", "Format output as: json, ndjson.")
	fl.StringArrayVar(&opts.placeholderFlags, "placeholder", nil, "Define a custom placeholder in <name>=<value> format, expanded from :<name>. Repeat the flag to define more than one.")

	cmd.MarkFlagsMutuallyExclusive("paginate", "input")
	cmd.MarkFlagsMutuallyExclusive("form", "field")
	cmd.MarkFlagsMutuallyExclusive("form", "raw-field")
	cmd.MarkFlagsMutuallyExclusive("form", "input")

	return cmd
}
