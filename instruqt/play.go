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

// playType defines a custom type for play modes on Instruqt.
type playType string

// Constants representing different types of plays.
const (
	PlayTypeAll       playType = "ALL"       // Represents all play types.
	PlayTypeDeveloper playType = "DEVELOPER" // Represents developer-specific play types.
	PlayTypeNormal    playType = "NORMAL"    // Represents normal play types.
)

// playQuery represents the GraphQL query structure for retrieving play reports
// with specific filters like team slug, date range, and pagination.
type playQuery struct {
	PlayReports `graphql:"playReports(input: {teamSlug: $teamSlug, dateRangeFilter: {from: $from, to: $to}, pagination: {skip: $skip, take: $take}})"`
}

// PlayReports represents a collection of play reports retrieved from Instruqt.
type PlayReports struct {
	Items      []PlayReport // A list of play reports.
	TotalItems int          // The total number of play reports available.
}

// PlayReport represents the data structure for a single play report on Instruqt.
type PlayReport struct {
	Id    string // The unique identifier for the play report.
	Track struct {
		Id string // The unique identifier of the track associated with the play.
	}
	TrackInvite struct {
		Id string // The unique identifier of the track invite associated with the play.
	}
	User struct {
		Id string // The unique identifier of the user associated with the play.
	}

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

// GetPlays retrieves a list of play reports from Instruqt for the specified team,
// within a given date range, and using pagination parameters.
//
// Parameters:
//   - from: The start date of the date range filter.
//   - to: The end date of the date range filter.
//   - take: The number of play reports to retrieve in one call.
//   - skip: The number of play reports to skip before starting to retrieve.
//
// Returns:
//   - []PlayReport: A list of play reports that match the given criteria.
//   - int: The total number of play reports available for the given criteria.
//   - error: Any error encountered while retrieving the play reports.
func (c *Client) GetPlays(from time.Time, to time.Time, take int, skip int) (plays []PlayReport, totalItems int, err error) {
	variables := map[string]interface{}{
		"teamSlug": graphql.String(c.TeamSlug),
		"from":     from,
		"to":       to,
		"take":     graphql.Int(take),
		"skip":     graphql.Int(skip),
	}

	var q playQuery
	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return plays, 0, err
	}

	return q.PlayReports.Items, q.TotalItems, nil
}
