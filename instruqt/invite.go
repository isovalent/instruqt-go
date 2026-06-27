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

// inviteQuery represents the GraphQL query structure for retrieving a single
// track invite by its invite ID.
type inviteQuery struct {
	TrackInvite `graphql:"trackInvite(inviteID: $inviteId)"`
}

type inviteTracksQuery struct {
	TrackInvite trackInviteTracks `graphql:"trackInvite(inviteID: $inviteId)"`
}

// TrackInvite represents the data structure for an Instruqt track invite.
type TrackInvite struct {
	Id                                 string    // The unique identifier for the invite.
	Title                              string    // The internal title of the track invite.
	PublicTitle                        string    // The public title of the track invite.
	PublicDescription                  string    // The public description of the track invite.
	AccessSetting                      string    // The access setting for the invite.
	InviteLimit                        int       // The maximum number of claims allowed for the invite.
	InviteCount                        int       // The number of times the invite has been used.
	ClaimCount                         int       // The number of claims associated with the invite.
	ExpiresAt                          time.Time // The timestamp when the invite expires.
	StartsAt                           time.Time // The timestamp when the invite becomes available.
	Created                            time.Time // The timestamp when the invite was created.
	Last_Updated                       time.Time // The timestamp when the invite was last updated.
	AllowAnonymous                     bool      // Whether anonymous users can claim the invite.
	AllowedEmailAddresses              []string  // The email addresses allowed to claim the invite.
	AllowedEmailAddressesOnly          bool      // Whether only explicitly allowed email addresses can claim the invite.
	CurrentUserAllowed                 bool      // Whether the current API user can claim the invite.
	CurrentUserClaimed                 bool      // Whether the current API user has claimed the invite.
	Type                               string    // The invite type.
	Status                             string    // The invite status.
	DaysUntil                          int       // Number of days until the invite starts or expires.
	CanClaim                           bool      // Whether the invite can currently be claimed.
	EmailOwnershipConfirmationRequired bool      // Whether email ownership confirmation is required.
	RuntimeParameters                  struct {  // The runtime parameters associated with the invite.
		EnvironmentVariables []variable // Environment variables used during the invite session.
	}
	Claims []TrackInviteClaim // A list of claims associated with the track invite.
	Tracks []Track            `graphql:"-"` // A list of tracks associated with the invite, only queried with WithTracks().
}

type trackInviteTracks struct {
	Id     string
	Tracks []Track
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
//
//   - inviteId: The unique identifier of the track invite to retrieve.
//
//   - opts: Optional query modifiers, such as WithTracks.
//
// Returns:
//   - TrackInvite: The track invite details if found.
//   - error: Any error encountered while retrieving the invite.
func (c *Client) GetInvite(inviteId string, opts ...Option) (i TrackInvite, err error) {
	if inviteId == "" {
		return i, nil
	}

	variables := map[string]interface{}{
		"inviteId": graphql.String(inviteId),
	}

	var q inviteQuery
	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return i, err
	}

	options := &options{}
	for _, opt := range opts {
		opt(options)
	}
	if options.includeTracks {
		tracks, err := c.GetInviteTracks(inviteId)
		if err != nil {
			return i, err
		}
		q.TrackInvite.Tracks = tracks
	}

	return q.TrackInvite, nil
}

// GetInviteTracks retrieves the tracks associated with a track invite.
//
// Parameters:
//   - inviteId: The unique identifier of the track invite to retrieve tracks for.
//
// Returns:
//   - []Track: The tracks associated with the invite.
//   - error: Any error encountered while retrieving the invite tracks.
func (c *Client) GetInviteTracks(inviteId string) ([]Track, error) {
	if inviteId == "" {
		return nil, nil
	}

	var q inviteTracksQuery
	variables := map[string]interface{}{
		"inviteId": graphql.String(inviteId),
	}
	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return nil, err
	}

	return q.TrackInvite.Tracks, nil
}

// invitesQuery represents the GraphQL query structure for retrieving all track invites
// for a specific team.
type invitesQuery struct {
	TrackInvites []TrackInvite `graphql:"trackInvites(teamSlug: $teamSlug)"`
}

type invitesTracksQuery struct {
	TrackInvites []trackInviteTracks `graphql:"trackInvites(teamSlug: $teamSlug)"`
}

// GetInvites retrieves all track invites for the specified team slug from Instruqt.
//
// Returns:
//   - []TrackInvite: A list of track invites for the team.
//   - error: Any error encountered while retrieving the invites.
func (c *Client) GetInvites(opts ...Option) (i []TrackInvite, err error) {
	variables := map[string]interface{}{
		"teamSlug": graphql.String(c.TeamSlug),
	}

	var q invitesQuery
	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return i, err
	}

	options := &options{}
	for _, opt := range opts {
		opt(options)
	}
	if options.includeTracks {
		tracksByInvite, err := c.GetInvitesTracks()
		if err != nil {
			return i, err
		}
		for idx := range q.TrackInvites {
			q.TrackInvites[idx].Tracks = tracksByInvite[q.TrackInvites[idx].Id]
		}
	}

	return q.TrackInvites, nil
}

// GetInvitesTracks retrieves track lists for all track invites in the team.
//
// Returns:
//   - map[string][]Track: A map of invite ID to associated tracks.
//   - error: Any error encountered while retrieving invite tracks.
func (c *Client) GetInvitesTracks() (map[string][]Track, error) {
	var q invitesTracksQuery
	variables := map[string]interface{}{
		"teamSlug": graphql.String(c.TeamSlug),
	}
	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return nil, err
	}

	tracksByInvite := make(map[string][]Track, len(q.TrackInvites))
	for _, invite := range q.TrackInvites {
		tracksByInvite[invite.Id] = invite.Tracks
	}
	return tracksByInvite, nil
}
