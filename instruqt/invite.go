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

	"github.com/shurcooL/graphql"
)

// inviteQuery represents the GraphQL query structure for retrieving a single
// track invite by its invite ID.
type inviteQuery struct {
	TrackInvite `graphql:"trackInvite(inviteID: $inviteId)"`
}

// TrackInvite represents the data structure for an Instruqt track invite.
type TrackInvite struct {
	Id                string   // The unique identifier for the invite.
	PublicTitle       string   // The public title of the track invite.
	RuntimeParameters struct { // The runtime parameters associated with the invite.
		EnvironmentVariables []variable // Environment variables used during the invite session.
	}
	Claims []TrackInviteClaim // A list of claims associated with the track invite.
}

// TrackInviteClaim represents a claim made by a user for a specific track invite.
type TrackInviteClaim struct {
	Id   string   // The unique identifier of the claim.
	User struct { // Information about the user who made the claim.
		Id string // The unique identifier of the user.
	}
	ClaimedAt time.Time // The timestamp when the claim was made.
}

// variable represents an environment variable key-value pair.
type variable struct {
	Key   string // The key of the environment variable.
	Value string // The value of the environment variable.
}

// GetInvite retrieves a track invite from Instruqt using its unique invite ID.
//
// Parameters:
//   - inviteId: The unique identifier of the track invite to retrieve.
//
// Returns:
//   - TrackInvite: The track invite details if found.
//   - error: Any error encountered while retrieving the invite.
func (c *Client) GetInvite(inviteId string) (i TrackInvite, err error) {
	if inviteId == "" {
		return i, nil
	}

	var q inviteQuery
	variables := map[string]interface{}{
		"inviteId": graphql.String(inviteId),
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return i, err
	}

	return q.TrackInvite, nil
}

// invitesQuery represents the GraphQL query structure for retrieving all track invites
// for a specific team.
type invitesQuery struct {
	TrackInvites []TrackInvite `graphql:"trackInvites(teamSlug: $teamSlug)"`
}

// GetInvites retrieves all track invites for the specified team slug from Instruqt.
//
// Returns:
//   - []TrackInvite: A list of track invites for the team.
//   - error: Any error encountered while retrieving the invites.
func (c *Client) GetInvites() (i []TrackInvite, err error) {
	var q invitesQuery
	variables := map[string]interface{}{
		"teamSlug": graphql.String(c.TeamSlug),
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return i, err
	}

	return q.TrackInvites, nil
}
