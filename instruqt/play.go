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

	graphql "github.com/hasura/go-graphql-client"
)

// playType defines a custom type for play modes on Instruqt.
type PlayType string

// Constants representing different types of plays.
const (
	PlayTypeAll       PlayType = "ALL"       // Represents all play types.
	PlayTypeDeveloper PlayType = "DEVELOPER" // Represents developer-specific play types.
	PlayTypeNormal    PlayType = "NORMAL"    // Represents normal play types.
)

// playQuery represents the GraphQL query structure for retrieving play reports
// with specific filters like team slug, date range, and pagination.
type playQuery struct {
	PlayReports `graphql:"playReports(input: {teamSlug: $teamSlug, dateRangeFilter: {from: $from, to: $to}, trackIds: $trackIds, trackInviteIds: $trackInviteIds, landingPageIds: $landingPageIds, tags: $tags,  userIds: $userIds, pagination: {skip: $skip, take: $take}, playType: $playType, customParameterFilters: $customParameterFilters, ordering: {orderBy: $orderBy, direction: $orderDirection}})"`
}

// Play is the domain model of a user's journey through a track.
type Play struct {
	Id        string
	StartedAt time.Time
}

// PlayReports represents a collection of play reports retrieved from Instruqt.
type PlayReports struct {
	Items      []PlayReport // A list of play reports.
	TotalItems int          // The total number of play reports available.
}

// PlayReport represents the data structure for a single play report on Instruqt.
type PlayReport struct {
	Id    string       // The unique identifier for the play report.
	Track SandboxTrack // The track played.

	TrackInvite TrackInvite // The optional Track invite associated to the play.

	User User // The user that played the play.

	CompletionPercent   float64   // The percentage of the play that has been completed.
	TotalChallenges     int       // The total number of challenges in the play.
	CompletedChallenges int       // The number of challenges completed by the user in the play.
	TimeSpent           int       // The total time spent on the play, in seconds.
	StoppedReason       string    // The reason why the play was stopped (if applicable).
	Mode                string    // The mode of the play (e.g., NORMAL, DEVELOPER).
	StartedAt           time.Time // The time when the play started.

	Activity []struct { // A list of activities performed during the play.
		Time    time.Time // The time when the activity occurred.
		Message string    // A message describing the activity.
	}

	PlayReview struct { // The review details for the play, if available.
		Id      string // The unique identifier of the play review.
		Score   int    // The score given in the play review.
		Content string // The content of the play review.
	}

	CustomParameters []struct { // Custom parameters associated with the play.
		Key   string // The key of the custom parameter.
		Value string // The value of the custom parameter.
	}
}

// playItemQuery represents the GraphQL query structure for retrieving a single play report
type playItemQuery struct {
	PlayReportItem PlayReport `graphql:"playReportItem(playID: $playID, input: {teamSlug: $teamSlug, playType: $playType})"`
}

// GetPlays retrieves a list of play reports from Instruqt for the specified team,
// within a given date range, and using pagination parameters.
//
// Parameters:
//   - from: The start date of the date range filter.
//   - to: The end date of the date range filter.
//   - take: The number of play reports to retrieve in one call.
//   - skip: The number of play reports to skip before starting to retrieve.
//   - opts: A variadic number of Option to configure the query.
//
// Returns:
//   - []PlayReport: A list of play reports that match the given criteria.
//   - int: The total number of play reports available for the given criteria.
//   - error: Any error encountered while retrieving the play reports.
func (c *Client) GetPlays(from time.Time, to time.Time, take int, skip int, opts ...Option) ([]PlayReport, int, error) {
	// Initialize the filter with default values
	filters := &options{
		trackIDs:               []string{},
		trackInviteIDs:         []string{},
		landingPageIDs:         []string{},
		tags:                   []string{},
		userIDs:                []string{},
		playType:               PlayTypeAll, // Default PlayType
		customParameterFilters: []CustomParameterFilter{},
		ordering: &Ordering{
			OrderBy:   OrderByCompletionPercent,
			Direction: DirectionDesc,
		},
	}

	// Apply each option to modify the filter
	for _, opt := range opts {
		opt(filters)
	}

	// Convert Go types to GraphQL types
	trackIds := make([]graphql.String, len(filters.trackIDs))
	for i, id := range filters.trackIDs {
		trackIds[i] = graphql.String(id)
	}

	trackInviteIds := make([]graphql.String, len(filters.trackInviteIDs))
	for i, id := range filters.trackInviteIDs {
		trackInviteIds[i] = graphql.String(id)
	}

	landingPageIds := make([]graphql.String, len(filters.landingPageIDs))
	for i, id := range filters.landingPageIDs {
		landingPageIds[i] = graphql.String(id)
	}

	tags := make([]graphql.String, len(filters.tags))
	for i, tag := range filters.tags {
		tags[i] = graphql.String(tag)
	}

	userIds := make([]graphql.String, len(filters.userIDs))
	for i, id := range filters.userIDs {
		userIds[i] = graphql.String(id)
	}

	customParameterFilters := make([]CustomParameterFilter, len(filters.customParameterFilters))
	for i, filter := range filters.customParameterFilters {
		customParameterFilters[i] = CustomParameterFilter{
			Key:   filter.Key,
			Value: filter.Value,
		}
	}

	// Prepare the variables map for the GraphQL query
	variables := map[string]interface{}{
		"teamSlug":               graphql.String(c.TeamSlug),
		"from":                   from,
		"to":                     to,
		"trackIds":               trackIds,
		"trackInviteIds":         trackInviteIds,
		"landingPageIds":         landingPageIds,
		"tags":                   tags,
		"userIds":                userIds,
		"take":                   graphql.Int(take),
		"skip":                   graphql.Int(skip),
		"playType":               filters.playType,
		"customParameterFilters": customParameterFilters,
		"orderBy":                graphql.String(filters.ordering.OrderBy),
		"orderDirection":         filters.ordering.Direction,
	}

	var q playQuery
	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return nil, 0, fmt.Errorf("GraphQL query failed: %w", err)
	}

	return q.PlayReports.Items, q.PlayReports.TotalItems, nil
}

func (c *Client) GetPlayReportItem(playId string, opts ...Option) (*PlayReport, error) {
	// Initialize the filter with default values
	filters := &options{
		playType: PlayTypeAll, // Default PlayType
	}

	// Apply each option to modify the filter
	for _, opt := range opts {
		opt(filters)
	}

	// Prepare the variables map for the GraphQL query
	variables := map[string]interface{}{
		"playID":   graphql.String(playId),
		"teamSlug": graphql.String(c.TeamSlug),
		"playType": filters.playType,
	}

	var q playItemQuery
	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}

	if filters.includeChallenges {
		challenges, err := c.GetChallenges(q.PlayReportItem.Track.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch challenges for track: %v", err)
		}
		q.PlayReportItem.Track.Challenges = challenges
	}

	return &q.PlayReportItem, nil
}
