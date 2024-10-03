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
	"fmt"
	"time"

	"github.com/shurcooL/graphql"
)

// baseReview represents the fundamental fields of a review.
// It includes the review ID, score, content, and timestamps for creation and updates.
type baseReview struct {
	Id         string    `json:"-"`          // The review ID.
	Score      int       `json:"score"`      // The score given in the review.
	Content    string    `json:"content"`    // The content of the review.
	Created_At time.Time `json:"created_at"` // The timestamp when the review was created.
	Updated_At time.Time `json:"updated_at"` // The timestamp when the review was last updated.
}

// Review represents a review for an Instruqt track.
type Review struct {
	baseReview
	Play *Play
}

// GetReviewOption defines a functional option for configuring GetReview.
// It allows modifying the behavior of GetReview, such as including additional fields.
type GetReviewOption func(*reviewOptions)

// reviewOptions holds configuration options for GetReview.
// Currently, it supports whether to include the 'play' field in the query.
type reviewOptions struct {
	includePlay bool // Determines if the 'play' field should be included in the query.
}

// WithPlay is a functional option that configures GetReview to include the 'play' field in the query.
// Usage: GetReview("reviewID", WithPlay())
func WithPlay() GetReviewOption {
	return func(opts *reviewOptions) {
		opts.includePlay = true
	}
}

// GetReview retrieves a single review by its unique identifier.
// It accepts optional functional options to include additional fields like 'play'.
//
// Parameters:
// - id (string): The unique identifier of the review.
// - opts (...Option): Variadic functional options to modify the query behavior.
//
// Returns:
// - *Review: A pointer to the retrieved Review. Includes Play if specified.
// - error: An error object if the query fails or the review is not found.
func (c *Client) GetReview(id string, opts ...GetReviewOption) (*Review, error) {
	// Initialize default options.
	options := &reviewOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Prepare GraphQL variables.
	variables := map[string]interface{}{
		"id": graphql.ID(id),
	}

	if options.includePlay {
		// Define the extended query struct with 'play'.
		var q struct {
			TrackReview Review `graphql:"trackReview(reviewID: $id)"`
		}

		// Execute the query.
		if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
			return nil, fmt.Errorf("GraphQL query with play failed: %w", err)
		}

		// Return the fetched Review, which includes Play.
		return &q.TrackReview, nil
	}

	// Define the base query struct without 'play'.
	var q struct {
		TrackReview baseReview `graphql:"trackReview(reviewID: $id)"`
	}

	// Execute the query.
	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}

	// Construct the Review without Play.
	review := Review{
		baseReview: q.TrackReview,
		Play:       nil, // Play is not included.
	}

	return &review, nil
}
