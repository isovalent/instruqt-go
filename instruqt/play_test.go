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

	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPlays_NoFilters(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		TeamSlug:      "isovalent",
		Context:       context.Background(),
	}

	mockResponse := playQuery{
		PlayReports: PlayReports{
			Items:      []PlayReport{{Id: "play-1"}},
			TotalItems: 1,
		},
	}

	// Mock the query method for the GraphQL client
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		vars := args.Get(2).(map[string]interface{})

		q := args.Get(1).(*playQuery)
		*q = mockResponse

		// Check default play type is "ALL"
		assert.Equal(t, PlayTypeAll, vars["playType"])

		// Ensure no filters are passed
		assert.Empty(t, vars["trackIds"])
		assert.Empty(t, vars["trackInviteIds"])
	}).Return(nil).Once()

	// Execute with no filters
	from := time.Now().AddDate(0, 0, -30) // 30 days ago
	to := time.Now()
	plays, totalItems, err := client.GetPlays(from, to, 10, 0, nil)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, 1, totalItems)
	assert.Equal(t, "play-1", plays[0].Id)

	// Ensure the mock expectations are met
	mockClient.AssertExpectations(t)
}

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

	// Mock the query method for the GraphQL client
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		vars := args.Get(2).(map[string]interface{})

		q := args.Get(1).(*playQuery)
		*q = mockResponse

		// Check that the play type and filters are set correctly
		assert.Equal(t, PlayTypeDeveloper, vars["playType"])
		assert.Equal(t, []graphql.String{graphql.String("track-1")}, vars["trackIds"])
	}).Return(nil).Once()

	// Define filters to pass into the method
	filters := &PlayReportFilter{
		TrackIDs: []string{"track-1"},
		PlayType: PlayTypeDeveloper,
	}

	// Execute with filters
	from := time.Now().AddDate(0, 0, -30) // 30 days ago
	to := time.Now()
	plays, totalItems, err := client.GetPlays(from, to, 10, 0, filters)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, 1, totalItems)
	assert.Equal(t, "play-2", plays[0].Id)

	// Ensure the mock expectations are met
	mockClient.AssertExpectations(t)
}

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

	// Mock the query method for the GraphQL client
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		vars := args.Get(2).(map[string]interface{})

		q := args.Get(1).(*playQuery)
		*q = mockResponse

		// Check default play type is applied
		assert.Equal(t, PlayTypeAll, vars["playType"])

		// Ensure that track filters are applied correctly
		assert.Equal(t, []graphql.String{graphql.String("track-2")}, vars["trackIds"])
		assert.Empty(t, vars["trackInviteIds"])
	}).Return(nil).Once()

	// Define filters (only partial filters)
	filters := &PlayReportFilter{
		TrackIDs: []string{"track-2"},
	}

	// Execute with partial filters
	from := time.Now().AddDate(0, 0, -30) // 30 days ago
	to := time.Now()
	plays, totalItems, err := client.GetPlays(from, to, 10, 0, filters)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, 1, totalItems)
	assert.Equal(t, "play-3", plays[0].Id)

	// Ensure the mock expectations are met
	mockClient.AssertExpectations(t)
}

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
	mockClient.On("Query", mock.Anything, &playQuery{}, mock.Anything).Return(errors.New("graphql error"))

	// Call the method
	plays, totalItems, err := client.GetPlays(from, to, take, skip, nil)

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, plays)
	assert.Equal(t, 0, totalItems)
	assert.Contains(t, err.Error(), "graphql error")
	mockClient.AssertExpectations(t)
}
