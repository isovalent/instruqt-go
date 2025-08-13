// Copyright 2024 Cisco Systems, Inc. and its affiliates

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package instruqt

import (
	"context"
	"log"
	"net/http"
	"os"

	graphql "github.com/hasura/go-graphql-client"

	loghttp "github.com/motemen/go-loghttp"
)

// GraphQLClient is an interface that defines the methods for interacting with
// a GraphQL API, including querying and mutating data.
type GraphQLClient interface {
	Query(context.Context, any, map[string]any, ...graphql.Option) error
	Mutate(context.Context, any, map[string]any, ...graphql.Option) error
}

// Client represents the Instruqt API client, which provides methods to
// interact with the Instruqt platform. It includes a GraphQL client, logging capabilities,
// and the team slug to identify which team's data to interact with.
type Client struct {
	GraphQLClient GraphQLClient   // The GraphQL client used to execute queries and mutations.
	InfoLogger    *log.Logger     // Logger for informational messages.
	DebugLogger   *log.Logger     // Logger for debug messages.
	TeamSlug      string          // The slug identifier for the team within Instruqt.
	Context       context.Context // Default context for API requests
}

// NewClient creates a new instance of the Instruqt API client. It initializes
// the GraphQL client with the provided API token and team slug.
//
// Parameters:
//   - token: The API token used for authentication with the Instruqt GraphQL API.
//   - teamSlug: The slug identifier for the team.
//
// Returns:
//   - A pointer to the newly created Client instance.
func NewClient(token string, teamSlug string) *Client {
	client := &Client{
		InfoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		DebugLogger: log.New(os.Stdout, "DEBUG:", log.Ldate|log.Ltime),
		TeamSlug:    teamSlug,
		Context:     context.Background(), // Default context
	}

	httpClient := &http.Client{}
	httpClient.Transport = &BearerTokenRoundTripper{
		Transport: &loghttp.Transport{
			Transport: httpClient.Transport,
			/*
				LogRequest: func(req *http.Request) {
					b, _ := httputil.DumpRequestOut(req, true)
					client.DebugLogger.Printf("out body: %s", string(b))
				},
				LogResponse: func(resp *http.Response) {
					b, _ := httputil.DumpResponse(resp, true)
					client.DebugLogger.Printf("in body: %s", string(b))
				},
			*/
		},
		Token: token,
	}

	client.GraphQLClient = graphql.NewClient("https://play.instruqt.com/graphql", httpClient)
	return client
}

// WithContext creates a copy of the Client with a new context.
// This can be used to set specific timeouts or deadlines for API calls.
func (c *Client) WithContext(ctx context.Context) *Client {
	// Create a new Client instance with the same properties but a different context.
	return &Client{
		GraphQLClient: c.GraphQLClient,
		InfoLogger:    c.InfoLogger,
		TeamSlug:      c.TeamSlug,
		Context:       ctx,
	}
}

// BearerTokenRoundTripper is a custom HTTP RoundTripper that adds a Bearer token
// for authorization in the HTTP request headers.
type BearerTokenRoundTripper struct {
	Transport http.RoundTripper // The underlying transport to use for HTTP requests.
	Token     string            // The Bearer token for authorization.
}

// RoundTrip executes a single HTTP transaction, adding the Authorization header
// with the Bearer token to the request before forwarding it to the underlying transport.
//
// Parameters:
//   - req: The HTTP request to be sent.
//
// Returns:
//   - An HTTP response and any error encountered while making the request.
func (rt *BearerTokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+rt.Token)
	return rt.Transport.RoundTrip(req)
}
