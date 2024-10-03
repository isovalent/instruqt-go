// Copyright 2024 Cisco Systems, Inc. and its affiliates
//
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

	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGetPlays_WithFilters tests the GetPlays function with specific filters applied.
func TestGetPlays_WithFilters(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		TeamSlug:      "isovalent",
		Context:       context.Background(),
	}

	mockResponse := playQuery{
		PlayReports: PlayReports{
			Items:      []PlayReport{{Id: "play-2"}},
			TotalItems: 1,
		},
	}

	// Mock the Query method for the GraphQL client
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		vars := args.Get(2).(map[string]interface{})

		q := args.Get(1).(*playQuery)
		*q = mockResponse

		// Check that the play type and filters are set correctly
		assert.Equal(t, PlayTypeDeveloper, vars["playType"], "Expected playType to be PlayTypeDeveloper")
		assert.Equal(t, []graphql.String{graphql.String("track-1")}, vars["trackIds"], "Expected trackIds to contain 'track-1'")
	}).Return(nil).Once()

	// Define filters using Functional Options
	options := []Option{
		WithTrackIDs("track-1"),
		WithPlayType(PlayTypeDeveloper),
	}

	// Execute GetPlays with filters
	from := time.Now().AddDate(0, 0, -30) // 30 days ago
	to := time.Now()
	take := 10
	skip := 0
	plays, totalItems, err := client.GetPlays(from, to, take, skip, options...)

	// Assert the results
	assert.NoError(t, err, "Expected no error from GetPlays")
	assert.Equal(t, 1, totalItems, "Expected TotalItems to be 1")
	assert.Equal(t, "play-2", plays[0].Id, "Expected Play ID to be 'play-2'")

	// Ensure the mock expectations are met
	mockClient.AssertExpectations(t)
}

// TestGetPlays_WithPartialFilters tests the GetPlays function with partial filters applied.
func TestGetPlays_WithPartialFilters(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		TeamSlug:      "isovalent",
		Context:       context.Background(),
	}

	mockResponse := playQuery{
		PlayReports: PlayReports{
			Items:      []PlayReport{{Id: "play-3"}},
			TotalItems: 1,
		},
	}

	// Mock the Query method for the GraphQL client
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		vars := args.Get(2).(map[string]interface{})

		q := args.Get(1).(*playQuery)
		*q = mockResponse

		// Check default play type is applied
		assert.Equal(t, PlayTypeAll, vars["playType"], "Expected playType to be PlayTypeAll by default")

		// Ensure that track filters are applied correctly
		assert.Equal(t, []graphql.String{graphql.String("track-2")}, vars["trackIds"], "Expected trackIds to contain 'track-2'")
		assert.Empty(t, vars["trackInviteIds"], "Expected trackInviteIds to be empty")
	}).Return(nil).Once()

	// Define partial filters using Functional Options
	options := []Option{
		WithTrackIDs("track-2"),
	}

	// Execute GetPlays with partial filters
	from := time.Now().AddDate(0, 0, -30) // 30 days ago
	to := time.Now()
	take := 10
	skip := 0
	plays, totalItems, err := client.GetPlays(from, to, take, skip, options...)

	// Assert the results
	assert.NoError(t, err, "Expected no error from GetPlays with partial filters")
	assert.Equal(t, 1, totalItems, "Expected TotalItems to be 1 with partial filters")
	assert.Equal(t, "play-3", plays[0].Id, "Expected Play ID to be 'play-3'")

	// Ensure the mock expectations are met
	mockClient.AssertExpectations(t)
}

// TestGetPlays_Error tests the GetPlays function when the GraphQL client returns an error.
func TestGetPlays_Error(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 31, 23, 59, 59, 0, time.UTC)
	take := 10
	skip := 0

	// Set up the mock expectation to return an error
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("graphql error")).Once()

	// Call the method without any options
	plays, totalItems, err := client.GetPlays(from, to, take, skip)

	// Assertions
	assert.Error(t, err, "Expected an error from GetPlays")
	assert.Empty(t, plays, "Expected plays to be empty on error")
	assert.Equal(t, 0, totalItems, "Expected TotalItems to be 0 on error")
	assert.Contains(t, err.Error(), "graphql error", "Expected error message to contain 'graphql error'")

	// Ensure the mock expectations are met
	mockClient.AssertExpectations(t)
}
