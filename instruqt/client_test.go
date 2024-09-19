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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	token := "test-token"
	team := "my-amazing-team"
	client := NewClient(token, team)

	assert.NotNil(t, client.GraphQLClient, "GraphQLClient should not be nil")
	assert.IsType(t, &graphql.Client{}, client.GraphQLClient, "GraphQLClient should be of type *graphql.Client")
	assert.Equal(t, context.Background(), client.Context) // Check the default context

	// Create a request to test the custom RoundTripper
	req := httptest.NewRequest("GET", "http://example.com", nil)

	// Manually create the RoundTripper to ensure it uses the token
	rt := &BearerTokenRoundTripper{
		Transport: http.DefaultTransport,
		Token:     token,
	}

	// Execute the RoundTrip to verify the Authorization header
	_, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}

	// Check the Authorization header
	authHeader := req.Header.Get("Authorization")
	expectedAuthHeader := "Bearer " + token
	assert.Equal(t, expectedAuthHeader, authHeader, "Authorization header should be correctly set")
}

func TestClientWithContext(t *testing.T) {
	token := "test-token"
	teamSlug := "test-team"
	client := NewClient(token, teamSlug)

	// Create a new context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	clientWithCtx := client.WithContext(ctx)

	// Check that the new client has the updated context
	assert.NotNil(t, clientWithCtx)
	assert.Equal(t, ctx, clientWithCtx.Context)

	// Ensure the original client context remains unchanged
	assert.Equal(t, context.Background(), client.Context)
}

func TestGraphQLClientQueryWithContext(t *testing.T) {
	// Set up mock GraphQL client
	mockGraphQLClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockGraphQLClient,
		Context:       context.Background(),
	}

	// Define a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	clientWithCtx := client.WithContext(ctx)

	query := userInfoQuery{}
	variables := map[string]interface{}{
		"userID":   graphql.String("user-id"),
		"teamSlug": graphql.String(client.TeamSlug),
	}

	// Mock the expected behavior for the Query method
	mockGraphQLClient.On("Query", ctx, &query, variables).Return(nil)

	// Call GetUserInfo with the new client that has a custom context
	_, err := clientWithCtx.GetUserInfo("user-id")
	assert.NoError(t, err)

	// Verify that the Query method was called with the correct context
	mockGraphQLClient.AssertCalled(t, "Query", ctx, &query, variables)
}

func TestBearerTokenRoundTripper(t *testing.T) {
	token := "test-token"
	mockTransport := &mockRoundTripper{}
	rt := &BearerTokenRoundTripper{
		Transport: mockTransport,
		Token:     token,
	}

	req := httptest.NewRequest("GET", "http://example.com", nil)
	_, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	authHeader := req.Header.Get("Authorization")
	expectedAuthHeader := "Bearer " + token
	assert.Equal(t, expectedAuthHeader, authHeader, "Authorization header should be correctly set")

	// Verify that the request was passed to the mock transport
	assert.True(t, mockTransport.called, "Expected RoundTrip to be called on the underlying transport")
}

// mockRoundTripper is a mock implementation of http.RoundTripper used to test the BearerTokenRoundTripper
type mockRoundTripper struct {
	called bool
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.called = true
	return &http.Response{
		StatusCode: http.StatusOK,
	}, nil
}
