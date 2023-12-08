package main

import (
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"os"
	"time"
)

var (
	clientID     string
	clientSecret string
	dexConnector string
	oidcEndpoint string
	scopes       = []string{oidc.ScopeOpenID, "email", "federated:id"}
	token        string
)

const (
	tokenFile = "/run/secrets/kubernetes.io/serviceaccount/token"

	tokenTypeIdToken = "urn:ietf:params:oauth:token-type:id_token"
)

func main() {
	ctx := context.Background()
	// read Kubernetes Service Account token
	//tokenBytes, err := os.ReadFile(tokenFile)
	//if err != nil {
	//	log.Fatalf("could not read token file: %w", err)
	//	os.Exit(1)
	//}

	resp, err := ExchangeToken(
		ctx,
		oidcEndpoint,
		&TokenExchangeRequest{
			DexConnector:     dexConnector,
			Scope:            scopes,
			SubjectToken:     token,
			SubjectTokenType: tokenTypeIdToken,
		},
		ClientAuthentication{
			AuthStyle:    oauth2.AuthStyleInHeader,
			ClientID:     clientID,
			ClientSecret: clientSecret,
		},
		nil,
	)
	if err != nil {
		log.Fatalf("token exchange failed: %w", err)
		os.Exit(1)
	}

	tokenExpiration := time.Now().Local().Add(time.Duration(resp.ExpiresIn) * time.Second)

	fmt.Fprintln(os.Stderr, FormatExecCredential(resp.AccessToken, tokenExpiration))
	fmt.Println(FormatExecCredential(resp.AccessToken, tokenExpiration))
}
