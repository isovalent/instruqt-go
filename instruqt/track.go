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

// trackQuery represents the GraphQL query structure for retrieving a single
// track by its track ID.
type trackQuery struct {
	Track `graphql:"track(trackID: $trackId)"`
}

// userTrackQueryWithChallenges represents the GraphQL query structure for retrieving
// a user's specific track along with its challenges by track ID, user ID, and organization slug.
type userTrackQueryWithChallenges struct {
	Track SandboxTrack `graphql:"track(trackID: $trackId, userID: $userId, organizationSlug: $organizationSlug)"`
}

// trackQueryBySlug represents the GraphQL query structure for retrieving a single
// track by its slug and team slug.
type trackQueryBySlug struct {
	Track `graphql:"track(trackSlug: $trackSlug, teamSlug: $teamSlug)"`
}

// tracksQuery represents the GraphQL query structure for retrieving all tracks
// associated with a specific organization slug.
type tracksQuery struct {
	Tracks []Track `graphql:"tracks(organizationSlug: $organizationSlug)"`
}

// Track represents the data structure for an Instruqt track.
type Track struct {
	Slug        string   // The slug identifier for the track.
	Id          string   // The unique identifier for the track.
	Icon        string   // The icon associated with the track.
	Title       string   // The title of the track.
	Description string   // The description of the track.
	Teaser      string   // A teaser or short description of the track.
	Level       string   // The difficulty level of the track.
	Embed_Token string   // The token used for embedding the track.
	Statistics  struct { // Statistics about the track.
		Average_review_score float32 // The average review score of the track.
	}
	TrackTags    []TrackTag // A list of tags associated with the track.
	TrackReviews struct {   // A collection of reviews for the track.
		TotalCount int
		Nodes      []Review
	}
}

// TrackTag represents a tag associated with an Instruqt track.
type TrackTag struct {
	Value string // The value of the tag.
}

// SandboxTrack represents a track in a sandbox environment, including its details
// and associated challenges.
type SandboxTrack struct {
	Id          string   // The unique identifier for the sandbox track.
	Slug        string   // The slug identifier for the sandbox track.
	Icon        string   // The icon associated with the sandbox track.
	Title       string   // The title of the sandbox track.
	Description string   // The description of the sandbox track.
	Teaser      string   // A teaser or short description of the sandbox track.
	Level       string   // The difficulty level of the sandbox track.
	Embed_Token string   // The token used for embedding the sandbox track.
	Statistics  struct { // Statistics about the sandbox track.
		Average_review_score float32 // The average review score of the sandbox track.
	}
	TrackTags []struct { // A list of tags associated with the sandbox track.
		Value string
	}
	TrackReviews struct { // A collection of reviews for the sandbox track.
		TotalCount int
		Nodes      []Review
	}
	Challenges  []Challenge // A list of challenges associated with the sandbox track.
	Status      string      // The current status of the sandbox track.
	Started     time.Time   // The timestamp when the sandbox track was started.
	Completed   time.Time   // The timestamp when the sandbox track was completed.
	Participant struct {    // Information about the participant of the sandbox track.
		Id string
	}
}

// Review represents a review for an Instruqt track.
type Review struct {
	Score      int       `json:"score"`      // The score given in the review.
	Content    string    `json:"content"`    // The content of the review.
	Created_At time.Time `json:"created_at"` // The timestamp when the review was created.
	Updated_At time.Time `json:"updated_at"` // The timestamp when the review was last updated.
}

// GetTrackById retrieves a track from Instruqt using its unique track ID.
//
// Parameters:
//   - trackId: The unique identifier of the track to retrieve.
//
// Returns:
//   - Track: The track details if found.
//   - error: Any error encountered while retrieving the track.
func (c *Client) GetTrackById(trackId string) (t Track, err error) {
	if trackId == "" {
		return t, nil
	}

	var q trackQuery
	variables := map[string]interface{}{
		"trackId": graphql.String(trackId),
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return t, err
	}

	return q.Track, nil
}

// GetUserTrackById retrieves a track for a specific user, including its challenges,
// using the user's ID and the track's ID.
//
// Parameters:
//   - userId: The unique identifier of the user.
//   - trackId: The unique identifier of the track.
//
// Returns:
//   - SandboxTrack: The track details with challenges if found.
//   - error: Any error encountered while retrieving the track.
func (c *Client) GetUserTrackById(userId string, trackId string) (t SandboxTrack, err error) {
	if trackId == "" {
		return t, nil
	}

	var q userTrackQueryWithChallenges
	variables := map[string]interface{}{
		"trackId":          graphql.String(trackId),
		"userId":           graphql.String(userId),
		"organizationSlug": graphql.String(c.TeamSlug),
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return t, err
	}

	return q.Track, nil
}

// GetTrackBySlug retrieves a track from Instruqt using its slug and team slug.
//
// Parameters:
//   - trackSlug: The slug identifier of the track to retrieve.
//
// Returns:
//   - Track: The track details if found.
//   - error: Any error encountered while retrieving the track.
func (c *Client) GetTrackBySlug(trackSlug string) (t Track, err error) {
	if trackSlug == "" {
		return t, nil
	}

	var q trackQueryBySlug
	variables := map[string]interface{}{
		"trackSlug": graphql.String(trackSlug),
		"teamSlug":  graphql.String(c.TeamSlug),
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return t, err
	}

	return q.Track, nil
}

// GetTrackUnlockedChallenge retrieves the first unlocked challenge for a specific
// user's track.
//
// Parameters:
//   - userId: The unique identifier of the user.
//   - trackId: The unique identifier of the track.
//
// Returns:
//   - Challenge: The first unlocked challenge found.
//   - error: Any error encountered while retrieving the challenge.
func (c *Client) GetTrackUnlockedChallenge(userId string, trackId string) (challenge Challenge, err error) {
	track, err := c.GetUserTrackById(userId, trackId)
	if err != nil {
		return challenge, fmt.Errorf("[instruqt.GetTrackUnlockedChallenge] failed to get user track: %v", err)
	}

	for _, chllg := range track.Challenges {
		switch chllg.Status {
		case "unlocked", "creating", "created", "started":
			return chllg, nil
		}
	}

	return
}

// GetTracks retrieves all tracks associated with the client's team slug.
//
// Returns:
//   - []Track: A list of tracks for the team.
//   - error: Any error encountered while retrieving the tracks.
func (c *Client) GetTracks() (t []Track, err error) {
	var q tracksQuery
	variables := map[string]interface{}{
		"organizationSlug": graphql.String(c.TeamSlug),
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return t, err
	}

	return q.Tracks, nil
}

// GenerateOneTimePlayToken generates a one-time play token for a specific track.
//
// Parameters:
//   - trackId: The unique identifier of the track.
//
// Returns:
//   - string: The generated one-time play token.
//   - error: Any error encountered while generating the token.
func (c *Client) GenerateOneTimePlayToken(trackId string) (token string, err error) {
	var m struct {
		GenerateOneTimePlayToken string `graphql:"generateOneTimePlayToken(trackID: $trackID)"`
	}

	variables := map[string]any{
		"trackID": graphql.String(trackId),
	}

	if err := c.GraphQLClient.Mutate(c.Context, &m, variables); err != nil {
		return "", err
	}

	return m.GenerateOneTimePlayToken, nil
}
