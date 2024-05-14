package instruqt

import (
	"context"
	"time"

	"github.com/shurcooL/graphql"
)

type playType string

const (
	PlayTypeAll       playType = "ALL"
	PlayTypeDeveloper playType = "DEVELOPER"
	PlayTypeNormal    playType = "NORMAL"
)

type playQuery struct {
	PlayReports `graphql:"playReports(input: {teamSlug: $teamSlug, dateRangeFilter: {from: $from, to: $to}, pagination: {skip: $skip, take: $take}})"`
}

type PlayReports struct {
	Items      []PlayReport
	TotalItems int
}

type PlayReport struct {
	Id    string
	Track struct {
		Id string
	}
	TrackInvite struct {
		Id string
	}
	User struct {
		Id string
	}

	CompletionPercent   float64
	TotalChallenges     int
	CompletedChallenges int
	TimeSpent           int
	StoppedReason       string
	Mode                string
	StartedAt           time.Time

	Activity []struct {
		Time    time.Time
		Message string
	}

	PlayReview struct {
		Id      string
		Score   int
		Content string
	}

	CustomParameters []struct {
		Key   string
		Value string
	}
}

// TODO: playTypes
func (c *Client) GetPlays(from time.Time, to time.Time, take int, skip int) (plays []PlayReport, totalItems int, err error) {
	variables := map[string]interface{}{
		"teamSlug": graphql.String(teamSlug),
		"from":     from,
		"to":       to,
		"take":     graphql.Int(take),
		"skip":     graphql.Int(skip),
	}

	ctx := context.Background()

	var q playQuery
	if err := c.GraphQLClient.Query(ctx, &q, variables); err != nil {
		return plays, 0, err
	}

	return q.PlayReports.Items, q.TotalItems, nil
}
