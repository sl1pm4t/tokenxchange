/*
Copyright © 2023 Matt Morrison
*/
package cmd

import (
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/sl1pm4t/tokenxchange/credentials"
	"github.com/sl1pm4t/tokenxchange/exchange"
	"golang.org/x/oauth2"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	clientID     string
	clientSecret string
	dexConnector string
	oidcEndpoint string
	scopes       = []string{oidc.ScopeOpenID, "email", "federated:id"}
	tokenFile    string
)

const (
	ksaTokenFilePath = "/run/secrets/kubernetes.io/serviceaccount/token"
	tokenTypeIdToken = "urn:ietf:params:oauth:token-type:id_token"
)

func init() {
	rootCmd.Flags().StringVar(&clientID, "client-id", "", "OIDC Client ID")
	rootCmd.Flags().StringVar(&clientSecret, "client-secret", "", "OIDC Client Secret")
	rootCmd.Flags().StringVar(&dexConnector, "dex-connector", "", "Name of the Dex Connector")
	rootCmd.Flags().StringVar(&oidcEndpoint, "oidc-endpoint", "", "OIDC URL (e.g. Dex URL)")
	rootCmd.Flags().StringSliceVar(&scopes, "scopes", []string{oidc.ScopeOpenID, "email", "federated:id"}, "OIDC Client ID")
	rootCmd.Flags().StringVar(&tokenFile, "token-file", ksaTokenFilePath, "Full path to token file")

}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tokenxchange",
	Short: "Kubernetes token exchange auth helper",
	Long: `A kubernetes credential helper that exchanges an existing token
for a token signed by another OIDC issuer (e.g. Dex).`,

	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		// read Kubernetes Service Account token
		tokenBytes, err := os.ReadFile(tokenFile)
		if err != nil {
			log.Fatal(fmt.Errorf("could not read token file: %w", err))
		}

		resp, err := exchange.ExchangeToken(
			ctx,
			oidcEndpoint,
			&exchange.TokenExchangeRequest{
				DexConnector:     dexConnector,
				Scope:            scopes,
				SubjectToken:     string(tokenBytes),
				SubjectTokenType: tokenTypeIdToken,
			},
			exchange.ClientAuthentication{
				AuthStyle:    oauth2.AuthStyleInHeader,
				ClientID:     clientID,
				ClientSecret: clientSecret,
			},
			nil,
		)
		if err != nil {
			log.Fatal(fmt.Errorf("token exchange failed: %w", err))
		}

		tokenExpiration := time.Now().Local().Add(time.Duration(resp.ExpiresIn) * time.Second)

		fmt.Println(credentials.FormatExecCredential(resp.AccessToken, tokenExpiration))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
