package client

import (
	"context"
)

const (
// BBTokenServerUrl       = "https://awsextra.org/site/oauth2/access_token"
// BBApiServerUrl         = "https://api.awsextra.org/2.0"
// BBInternalApiServerUrl = "https://api.awsextra.org/internal"
// ApplicationJson        = "application/json"
// FormEncoded            = "application/x-www-form-urlencoded"
// Bearer                 = "Bearer"
// IdSeparator            = ":"
)

type Client struct {
	//Workspace    string
	//accessToken  string
	//clientId     string
	//clientSecret string
	//numRetries   int
	//retryDelay   int
	//httpClient   *http.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	c := &Client{
		//Workspace:    workspace,
		//accessToken:  accessToken,
		//clientId:     clientId,
		//clientSecret: clientSecret,
		//numRetries:   numRetries,
		//retryDelay:   retryDelay,
		//httpClient:   &http.Client{},
	}
	return c, nil
}
