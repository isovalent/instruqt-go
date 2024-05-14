package instruqt

import (
	"context"

	"github.com/shurcooL/graphql"
)

type challengeQuery struct {
	Challenge `graphql:"challenge(challengeID: $challengeId)"`
}

type Challenge struct {
	Id    string
	Title string
	Index int
	Track struct {
		Id string
	}
}

func (c *Client) GetChallenge(id string) (ch Challenge, err error) {
	if id == "" {
		return ch, nil
	}

	var q challengeQuery
	variables := map[string]interface{}{
		"challengeId": graphql.String(id),
	}

	ctx := context.Background()
	if err := c.GraphQLClient.Query(ctx, &q, variables); err != nil {
		return ch, err
	}

	return q.Challenge, nil
}
