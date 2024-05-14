package instruqt

import (
	"log"
	"net/http"

	"cloud.google.com/go/logging"
	"github.com/shurcooL/graphql"
)

const (
	teamSlug = "isovalent"
)

type Client struct {
	GraphQLClient *graphql.Client

	LogClient  *logging.Client
	InfoLogger *log.Logger
}

func NewClient(token string) *Client {
	httpClient := &http.Client{}
	httpClient.Transport = &BearerTokenRoundTripper{
		Transport: http.DefaultTransport,
		Token:     token,
	}
	return &Client{
		GraphQLClient: graphql.NewClient("https://play.instruqt.com/graphql", httpClient),
	}
}

type BearerTokenRoundTripper struct {
	Transport http.RoundTripper
	Token     string
}

func (rt *BearerTokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+rt.Token)
	return rt.Transport.RoundTrip(req)
}
