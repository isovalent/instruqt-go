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
	PlayReports `graphql:"playReports(input: {teamSlug: $teamSlug, dateRangeFilter: {from: $from, to: $to}, trackIds: $trackIds, trackInviteIds: $trackInviteIds, landingPageIds: $landingPageIds, tags: $tags,  userIds: $userIds, pagination: {skip: $skip, take: $take}, playType: $playType})"`
}

// PlayReportFilter defines the optional filters for fetching play reports.
type PlayReportFilter struct {
	TrackIDs       []string
	TrackInviteIDs []string
	LandingPageIDs []string
	Tags           []string
	UserIDs        []string
	PlayType       PlayType
}

// PlayReports represents a collection of play reports retrieved from Instruqt.
type PlayReports struct {
	Items      []PlayReport // A list of play reports.
	TotalItems int          // The total number of play reports available.
}

// PlayReport represents the data structure for a single play report on Instruqt.
type PlayReport struct {
	Id    string // The unique identifier for the play report.
	Track Track  // The track played.

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

// GetPlaysOption defines a function type that modifies PlayReportFilter
type GetPlaysOption func(*PlayReportFilter)

// WithTrackIDs sets the TrackIDs filter
func WithTrackIDs(ids ...string) GetPlaysOption {
	return func(f *PlayReportFilter) {
		f.TrackIDs = ids
	}
}

// WithTrackInviteIDs sets the TrackInviteIDs filter
func WithTrackInviteIDs(ids ...string) GetPlaysOption {
	return func(f *PlayReportFilter) {
		f.TrackInviteIDs = ids
	}
}

// WithTags sets the Tags filter
func WithTags(tags ...string) GetPlaysOption {
	return func(f *PlayReportFilter) {
		f.Tags = tags
	}
}

// WithUserIDs sets the UserIDs filter
func WithUserIDs(ids ...string) GetPlaysOption {
	return func(f *PlayReportFilter) {
		f.UserIDs = ids
	}
}

// WithPlayType sets the PlayType filter
func WithPlayType(pt PlayType) GetPlaysOption {
	return func(f *PlayReportFilter) {
		f.PlayType = pt
	}
}

// GetPlays retrieves a list of play reports from Instruqt for the specified team,
// within a given date range, and using pagination parameters.
//
// Parameters:
//   - from: The start date of the date range filter.
//   - to: The end date of the date range filter.
//   - take: The number of play reports to retrieve in one call.
//   - skip: The number of play reports to skip before starting to retrieve.
//   - opts: A variadic number of GetPlaysOption to configure the query.
//
// Returns:
//   - []PlayReport: A list of play reports that match the given criteria.
//   - int: The total number of play reports available for the given criteria.
//   - error: Any error encountered while retrieving the play reports.
//
// GetPlays retrieves play reports with optional filters, ordering, and custom parameters.
// It accepts a variadic number of GetPlaysOption to configure the query.
func (c *Client) GetPlays(from time.Time, to time.Time, take int, skip int, opts ...GetPlaysOption) ([]PlayReport, int, error) {
	// Initialize the filter with default values
	filters := &PlayReportFilter{
		TrackIDs:       []string{},
		TrackInviteIDs: []string{},
		LandingPageIDs: []string{},
		Tags:           []string{},
		UserIDs:        []string{},
		PlayType:       PlayTypeAll, // Default PlayType
	}

	// Apply each option to modify the filter
	for _, opt := range opts {
		opt(filters)
	}

	// Convert Go types to GraphQL types
	trackIds := make([]graphql.String, len(filters.TrackIDs))
	for i, id := range filters.TrackIDs {
		trackIds[i] = graphql.String(id)
	}

	trackInviteIds := make([]graphql.String, len(filters.TrackInviteIDs))
	for i, id := range filters.TrackInviteIDs {
		trackInviteIds[i] = graphql.String(id)
	}

	landingPageIds := make([]graphql.String, len(filters.LandingPageIDs))
	for i, id := range filters.LandingPageIDs {
		landingPageIds[i] = graphql.String(id)
	}

	tags := make([]graphql.String, len(filters.Tags))
	for i, tag := range filters.Tags {
		tags[i] = graphql.String(tag)
	}

	userIds := make([]graphql.String, len(filters.UserIDs))
	for i, id := range filters.UserIDs {
		userIds[i] = graphql.String(id)
	}

	// Prepare the variables map for the GraphQL query
	variables := map[string]interface{}{
		"teamSlug":       graphql.String(c.TeamSlug),
		"from":           from,
		"to":             to,
		"trackIds":       trackIds,
		"trackInviteIds": trackInviteIds,
		"landingPageIds": landingPageIds,
		"tags":           tags,
		"userIds":        userIds,
		"take":           graphql.Int(take),
		"skip":           graphql.Int(skip),
		"playType":       filters.PlayType,
	}

	var q playQuery
	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return nil, 0, fmt.Errorf("GraphQL query failed: %w", err)
	}

	return q.PlayReports.Items, q.PlayReports.TotalItems, nil
}
