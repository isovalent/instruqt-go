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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGetReview_Success tests the successful retrieval of a review.
func TestGetReview_Success(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		TeamSlug:      "isovalent",
		Context:       context.Background(),
	}

	// Define the expected response
	expectedReview := Review{
		Score:      5,
		Content:    "Excellent track! Learned a lot.",
		Created_At: time.Now().AddDate(0, -1, 0), // 1 month ago
		Updated_At: time.Now(),
	}

	// Mock the Query method
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		// Extract the query struct to populate
		q := args.Get(1).(*reviewQuery)
		// Assign the expected review to the query result
		q.TrackReview = expectedReview
	}).Return(nil).Once()

	// Execute GetReview
	review, err := client.GetReview("review-123")

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, review)
	assert.Equal(t, expectedReview.Score, review.Score)
	assert.Equal(t, expectedReview.Content, review.Content)
	assert.WithinDuration(t, expectedReview.Created_At, review.Created_At, time.Second)
	assert.WithinDuration(t, expectedReview.Updated_At, review.Updated_At, time.Second)

	// Ensure the mock expectations were met
	mockClient.AssertExpectations(t)
}

// TestGetReview_QueryError tests the GetReview function when the GraphQL query fails.
func TestGetReview_QueryError(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		TeamSlug:      "isovalent",
		Context:       context.Background(),
	}

	// Mock the Query method to return an error
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("network error")).Once()

	// Execute GetReview
	review, err := client.GetReview("review-123")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, review)
	assert.Contains(t, err.Error(), "GraphQL query failed: network error")

	// Ensure the mock expectations were met
	mockClient.AssertExpectations(t)
}

// TestGetReview_NotFound tests the GetReview function when the review is not found.
func TestGetReview_NotFound(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		TeamSlug:      "isovalent",
		Context:       context.Background(),
	}

	// Define the expected response for a not found review
	// Assuming that the API returns a Review with zero values if not found
	// Alternatively, the API might return an error; adjust accordingly based on actual API behavior

	expectedReview := Review{
		Score:      0,
		Content:    "",
		Created_At: time.Time{}, // Zero value
		Updated_At: time.Time{}, // Zero value
	}

	// Mock the Query method
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*reviewQuery)
		q.TrackReview = expectedReview
	}).Return(nil).Once()

	// Execute GetReview
	review, err := client.GetReview("non-existent-id")

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, review)
	assert.Equal(t, 0, review.Score)
	assert.Equal(t, "", review.Content)
	assert.Equal(t, time.Time{}, review.Created_At)
	assert.Equal(t, time.Time{}, review.Updated_At)

	// Ensure the mock expectations were met
	mockClient.AssertExpectations(t)
}
