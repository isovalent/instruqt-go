package instruqt

import (
	"context"

	"github.com/shurcooL/graphql"
)

type trackQuery struct {
	Track `graphql:"track(trackID: $trackId)"`
}

type Track struct {
	Id    string
	Slug  string
	Title string
}

func (c *Client) GetTrack(trackId string) (t Track, err error) {
	if trackId == "" {
		return t, nil
	}

	var q trackQuery
	variables := map[string]interface{}{
		"trackId": graphql.String(trackId),
	}

	ctx := context.Background()
	if err := c.GraphQLClient.Query(ctx, &q, variables); err != nil {
		return t, err
	}

	return q.Track, nil
}
