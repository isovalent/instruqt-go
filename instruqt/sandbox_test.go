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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSandboxVariable(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	playID := "sandbox-123"
	key := "MY_VAR"
	expectedValue := "value123"

	queryResult := sandboxVarQuery{
		GetSandboxVariable: SandboxVar{
			Key:   key,
			Value: expectedValue,
		},
	}

	mockClient.On("Query", mock.Anything, &sandboxVarQuery{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*sandboxVarQuery)
		*q = queryResult
	}).Return(nil)

	value, err := client.GetSandboxVariable(playID, key)

	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
	mockClient.AssertExpectations(t)
}

func TestGetSandboxVariable_Error(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	playID := "sandbox-123"
	key := "MY_VAR"

	mockClient.On("Query", mock.Anything, &sandboxVarQuery{}, mock.Anything).Return(errors.New("graphql error"))

	value, err := client.GetSandboxVariable(playID, key)

	assert.Error(t, err)
	assert.Empty(t, value)
	assert.Contains(t, err.Error(), "graphql error")
	mockClient.AssertExpectations(t)
}

func TestGetSandbox(t *testing.T) {
	// Create a mock GraphQL client
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		Context:       context.Background(),
		TeamSlug:      "isovalent", // Include a teamSlug in the client
	}

	// Define the expected sandbox response
	expectedSandbox := Sandbox{
		Last_Activity_At: time.Now(),
		State:            "active",
		Track: SandboxTrack{
			Id:    "track-123",
			Title: "Test Track",
		},
		Invite: TrackInvite{
			Id: "invite-123",
		},
	}

	// Set up the mock to return the expected sandbox data
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		query := args.Get(1).(*sandboxQuery)
		query.Sandbox = expectedSandbox
	}).Return(nil)

	// Call the GetSandbox method
	sandbox, err := client.GetSandbox("sandbox-123")

	// Validate the results
	assert.NoError(t, err)
	assert.Equal(t, expectedSandbox, sandbox)

	// Ensure the mock expectations are met
	mockClient.AssertExpectations(t)
}

func TestGetSandboxes(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	expectedSandboxes := []Sandbox{
		{
			Last_Activity_At: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			State:            "running",
			Track: SandboxTrack{
				Id:    "track-123",
				Title: "Track 1",
			},
			Invite: TrackInvite{
				Id:          "invite-123",
				PublicTitle: "Invite 1",
			},
		},
		{
			Last_Activity_At: time.Date(2023, 1, 2, 14, 30, 0, 0, time.UTC),
			State:            "completed",
			Track: SandboxTrack{
				Id:    "track-456",
				Title: "Track 2",
			},
			Invite: TrackInvite{
				Id:          "invite-456",
				PublicTitle: "Invite 2",
			},
		},
	}

	queryResult := sandboxesQuery{
		Sandboxes: struct {
			Nodes []Sandbox
		}{
			Nodes: expectedSandboxes,
		},
	}

	mockClient.On("Query", mock.Anything, &sandboxesQuery{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*sandboxesQuery)
		*q = queryResult
	}).Return(nil)

	sandboxes, err := client.GetSandboxes()

	assert.NoError(t, err)
	assert.Equal(t, expectedSandboxes, sandboxes)
	mockClient.AssertExpectations(t)
}

func TestGetSandboxes_Error(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	mockClient.On("Query", mock.Anything, &sandboxesQuery{}, mock.Anything).Return(errors.New("graphql error"))

	sandboxes, err := client.GetSandboxes()

	assert.Error(t, err)
	assert.Empty(t, sandboxes)
	assert.Contains(t, err.Error(), "graphql error")
	mockClient.AssertExpectations(t)
}
