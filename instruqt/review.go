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

// Review represents a review for an Instruqt track.
type Review struct {
	Score      int       `json:"score"`      // The score given in the review.
	Content    string    `json:"content"`    // The content of the review.
	Created_At time.Time `json:"created_at"` // The timestamp when the review was created.
	Updated_At time.Time `json:"updated_at"` // The timestamp when the review was last updated.
}

// reviewQuery represents the GraphQL query structure for fetching a single review.
type reviewQuery struct {
	TrackReview Review `graphql:"trackReview(reviewID: $id)"`
}

// GetReview retrieves a single review by its unique identifier.
//
// Parameters:
// - id (string): The unique identifier of the review.
//
// Returns:
// - *Review: A pointer to the retrieved Review.
// - error: An error object if the query fails or the review is not found.
func (c *Client) GetReview(id string) (*Review, error) {
	var q reviewQuery
	variables := map[string]interface{}{
		"id": graphql.ID(id),
	}

	// Execute the GraphQL query
	err := c.GraphQLClient.Query(c.Context, &q, variables)
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}

	return &q.TrackReview, nil
}
