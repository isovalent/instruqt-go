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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPlays(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	// Define input parameters
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 31, 23, 59, 59, 0, time.UTC)
	take := 10
	skip := 0

	// Define expected PlayReports
	expectedPlays := []PlayReport{
		{
			Id: "play-123",
			Track: struct{ Id string }{
				Id: "track-123",
			},
			CompletionPercent:   100,
			TotalChallenges:     5,
			CompletedChallenges: 5,
			TimeSpent:           120,
			Mode:                "NORMAL",
			StartedAt:           time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			Id: "play-456",
			Track: struct{ Id string }{
				Id: "track-456",
			},
			CompletionPercent:   75,
			TotalChallenges:     4,
			CompletedChallenges: 3,
			TimeSpent:           90,
			Mode:                "DEVELOPER",
			StartedAt:           time.Date(2023, 1, 20, 15, 30, 0, 0, time.UTC),
		},
	}

	queryResult := playQuery{
		PlayReports: PlayReports{
			Items:      expectedPlays,
			TotalItems: len(expectedPlays),
		},
	}

	// Set up the mock expectation
	mockClient.On("Query", mock.Anything, &playQuery{}, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*playQuery)
		*q = queryResult
	}).Return(nil)

	// Call the method
	plays, totalItems, err := client.GetPlays(from, to, take, skip)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedPlays, plays)
	assert.Equal(t, len(expectedPlays), totalItems)
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
	plays, totalItems, err := client.GetPlays(from, to, take, skip)

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, plays)
	assert.Equal(t, 0, totalItems)
	assert.Contains(t, err.Error(), "graphql error")
	mockClient.AssertExpectations(t)
}
