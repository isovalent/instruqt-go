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

// TestGetReview_Success tests the successful retrieval of a review without including 'play'.
func TestGetReview_Success(t *testing.T) {
	// Initialize the mock GraphQL client.
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		Context:       context.Background(),
	}

	// Define the expected baseReview response.
	expectedBaseReview := baseReview{
		Id:         "review123",
		Score:      5,
		Content:    "Excellent track! Learned a lot.",
		Created_At: time.Now().AddDate(0, -1, 0), // 1 month ago
		Updated_At: time.Now(),
	}

	// Mock the Query method to return the expected baseReview.
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		// Extract the query struct to populate.
		q := args.Get(1).(*struct {
			TrackReview baseReview `graphql:"trackReview(reviewID: $id)"`
		})
		// Assign the expected baseReview to the query result.
		q.TrackReview = expectedBaseReview
	}).Return(nil).Once()

	// Execute GetReview without including 'play'.
	review, err := client.GetReview("review123")

	// Assertions to ensure no error and correct data is returned.
	assert.NoError(t, err)
	assert.NotNil(t, review)
	assert.Equal(t, expectedBaseReview.Id, review.Id)
	assert.Equal(t, expectedBaseReview.Score, review.Score)
	assert.Equal(t, expectedBaseReview.Content, review.Content)
	assert.WithinDuration(t, expectedBaseReview.Created_At, review.Created_At, time.Second)
	assert.WithinDuration(t, expectedBaseReview.Updated_At, review.Updated_At, time.Second)
	assert.Nil(t, review.Play) // Play should be nil since it's not included.

	// Ensure the mock expectations were met.
	mockClient.AssertExpectations(t)
}

// TestGetReview_Success_WithPlay tests the successful retrieval of a review including 'play'.
func TestGetReview_Success_WithPlay(t *testing.T) {
	// Initialize the mock GraphQL client.
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		Context:       context.Background(),
	}

	// Define the expected Review response with Play.
	expectedReview := Review{
		baseReview: baseReview{
			Id:         "review123",
			Score:      5,
			Content:    "Excellent track! Learned a lot.",
			Created_At: time.Now().AddDate(0, -1, 0), // 1 month ago
			Updated_At: time.Now(),
		},
		Play: &Play{
			Id: "play456",
		},
	}

	// Mock the Query method to return the expected Review with Play.
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		// Extract the query struct to populate.
		q := args.Get(1).(*struct {
			TrackReview Review `graphql:"trackReview(reviewID: $id)"`
		})
		// Assign the expected Review to the query result.
		q.TrackReview = expectedReview
	}).Return(nil).Once()

	// Execute GetReview with the WithPlay option to include 'play'.
	review, err := client.GetReview("review123", WithPlay())

	// Assertions to ensure no error and correct data is returned.
	assert.NoError(t, err)
	assert.NotNil(t, review)
	assert.Equal(t, expectedReview.Id, review.Id)
	assert.Equal(t, expectedReview.Score, review.Score)
	assert.Equal(t, expectedReview.Content, review.Content)
	assert.WithinDuration(t, expectedReview.Created_At, review.Created_At, time.Second)
	assert.WithinDuration(t, expectedReview.Updated_At, review.Updated_At, time.Second)
	assert.NotNil(t, review.Play) // Play should be included.
	assert.Equal(t, expectedReview.Play.Id, review.Play.Id)

	// Ensure the mock expectations were met.
	mockClient.AssertExpectations(t)
}

// TestGetReview_QueryError tests the GetReview function when the GraphQL query fails.
func TestGetReview_QueryError(t *testing.T) {
	// Initialize the mock GraphQL client.
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		Context:       context.Background(),
	}

	// Mock the Query method to return an error.
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("network error")).Once()

	// Execute GetReview without including 'play'.
	review, err := client.GetReview("review-123")

	// Assertions to ensure an error is returned and review is nil.
	assert.Error(t, err)
	assert.Nil(t, review)
	assert.Contains(t, err.Error(), "GraphQL query failed: network error")

	// Ensure the mock expectations were met.
	mockClient.AssertExpectations(t)
}

// TestGetReview_NotFound tests the GetReview function when the review is not found.
func TestGetReview_NotFound(t *testing.T) {
	// Initialize the mock GraphQL client.
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
		Context:       context.Background(),
	}

	// Define the expected response for a not found review.
	// Assuming the API returns zero values for a non-existent review.
	expectedBaseReview := baseReview{
		Id:         "",
		Score:      0,
		Content:    "",
		Created_At: time.Time{}, // Zero value
		Updated_At: time.Time{}, // Zero value
	}

	// Mock the Query method to return the expected baseReview with zero values.
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		q := args.Get(1).(*struct {
			TrackReview baseReview `graphql:"trackReview(reviewID: $id)"`
		})
		q.TrackReview = expectedBaseReview
	}).Return(nil).Once()

	// Execute GetReview with a non-existent ID.
	review, err := client.GetReview("non-existent-id")

	// Assertions to ensure no error is returned and review contains zero values.
	assert.NoError(t, err)
	assert.NotNil(t, review)
	assert.Equal(t, expectedBaseReview.Id, review.Id)
	assert.Equal(t, expectedBaseReview.Score, review.Score)
	assert.Equal(t, expectedBaseReview.Content, review.Content)
	assert.Equal(t, expectedBaseReview.Created_At, review.Created_At)
	assert.Equal(t, expectedBaseReview.Updated_At, review.Updated_At)
	assert.Nil(t, review.Play) // Play should be nil since it's not included.

	// Ensure the mock expectations were met.
	mockClient.AssertExpectations(t)
}
