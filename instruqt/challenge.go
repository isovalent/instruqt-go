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
	"time"

	graphql "github.com/hasura/go-graphql-client"
)

// challengeQuery represents the GraphQL query structure for retrieving a single challenge
// by its challenge ID.
type challengeQuery struct {
	Challenge `graphql:"challenge(challengeID: $challengeId)"`
}

// userChallengeQuery represents the GraphQL query structure for retrieving a single challenge
// associated with a specific user by user ID and challenge ID.
type userChallengeQuery struct {
	Challenge `graphql:"challenge(userID: $userId, challengeID: $challengeId)"`
}

// Challenge represents the data structure for an Instruqt challenge.
type Challenge struct {
	Id     string `json:"id"`     // The unique identifier for the challenge.
	Slug   string `json:"slug"`   // The slug for the challenge, which is a human-readable identifier.
	Title  string `json:"title"`  // The title of the challenge.
	Teaser string `json:"teaser"` // The teaser of the challenge.
	Index  int    `json:"index"`  // The index of the challenge in the track.
	Status string `json:"status"` // The status of the challenge (e.g., "unlocked", "completed").
	Track  struct {
		Id string // The identifier for the track associated with the challenge.
	} `json:"-"`
	Attempts []struct {
		Message   string    `json:"message"`   // The message returned by the attempts.
		Timestamp time.Time `json:"timestamp"` // The timestamp of the attempt.
	} `json:"attempts"` // The attempts made on the challenge by the user.
	Assignment string `json:"assignment"` // The assignment details for the challenge.
}

// GetChallenge retrieves a challenge from Instruqt using its unique challenge ID.
//
// Parameters:
//   - id: The unique identifier of the challenge to retrieve.
//
// Returns:
//   - Challenge: The challenge details if found.
//   - error: Any error encountered while retrieving the challenge.
func (c *Client) GetChallenge(id string) (ch Challenge, err error) {
	if id == "" {
		return ch, nil
	}

	var q challengeQuery
	variables := map[string]interface{}{
		"challengeId": graphql.String(id),
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return ch, err
	}

	return q.Challenge, nil
}

// GetUserChallenge retrieves a challenge associated with a specific user from Instruqt
// using the user's ID and the challenge's ID.
//
// Parameters:
//   - userId: The unique identifier of the user.
//   - id: The unique identifier of the challenge.
//
// Returns:
//   - Challenge: The challenge details if found.
//   - error: Any error encountered while retrieving the challenge.
func (c *Client) GetUserChallenge(userId string, id string) (ch Challenge, err error) {
	if id == "" {
		return ch, nil
	}

	var q userChallengeQuery
	variables := map[string]interface{}{
		"challengeId": graphql.String(id),
		"userId":      graphql.String(userId),
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return ch, err
	}

	return q.Challenge, nil
}

// SkipToChallenge allows a user to skip to a specific challenge in a track on Instruqt.
//
// Parameters:
//   - userId: The unique identifier of the user.
//   - trackId: The unique identifier of the track.
//   - id: The unique identifier of the challenge to skip to.
//
// Returns:
//   - error: Any error encountered while performing the skip operation.
func (c *Client) SkipToChallenge(userId string, trackId string, id string) (err error) {
	var m struct {
		SkipToChallenge struct {
			Id     graphql.String
			Status graphql.String
		} `graphql:"skipToChallenge(trackID: $trackID, challengeID: $challengeID, userID: $userID)"`
	}

	variables := map[string]any{
		"trackID":     graphql.String(trackId),
		"challengeID": graphql.String(id),
		"userID":      graphql.String(userId),
	}

	if err := c.GraphQLClient.Mutate(c.Context, &m, variables); err != nil {
		return err
	}

	return nil
}
