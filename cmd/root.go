package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Martinvks/httptestrunner/types"
	"github.com/spf13/cobra"
)

type commonArguments struct {
	addIdHeader   bool
	commonHeaders types.Headers
	keyLogFile    string
	proto         int
	timeout       time.Duration
	target        *url.URL
}

var (
	headers    []string
	proto      string
	commonArgs commonArguments
)

func init() {
	rootCmd.PersistentFlags().BoolVar(
		&commonArgs.addIdHeader,
		"id-header",
		false,
		"add a header field with name \"x-id\" and a uuid v4 value. the value will be added to the output when using the \"multi\" command",
	)

	rootCmd.PersistentFlags().StringArrayVarP(
		&headers,
		"header",
		"H",
		[]string{},
		"common header fields added to each request. syntax similar to curl: -H \"x-extra-header: val\"",
	)

	rootCmd.PersistentFlags().StringVarP(
		&commonArgs.keyLogFile,
		"keylogfile",
		"k",
		"",
		"filename to log TLS master secrets",
	)

	rootCmd.PersistentFlags().DurationVarP(
		&commonArgs.timeout,
		"timeout",
		"t",
		10*time.Second,
		"timeout",
	)

	rootCmd.PersistentFlags().StringVarP(
		&proto,
		"protocol",
		"p",
		"h2",
		"specifies which protocol to use. Must be one of \"h2\" or \"h3\"",
	)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

var rootCmd = &cobra.Command{
	Use:   "httptestrunner",
	Short: "An HTTP client for sending (possibly malformed) HTTP/2 and HTTP/3 requests",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		for _, header := range headers {
			parts := strings.Split(header, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid header '%s', expected syntax: 'x-extra-header: val'", header)
			}
			commonArgs.commonHeaders = append(
				commonArgs.commonHeaders,
				types.Header{
					Name:  strings.TrimSpace(strings.ToLower(parts[0])),
					Value: strings.TrimSpace(parts[1]),
				})
		}

		switch proto {
		case "h2":
			commonArgs.proto = types.H2
		case "h3":
			commonArgs.proto = types.H3
		default:
			return fmt.Errorf("unknown protocol '%s'", proto)
		}

		target, err := url.Parse(args[0])
		if err != nil {
			return err
		}
		commonArgs.target = target

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
