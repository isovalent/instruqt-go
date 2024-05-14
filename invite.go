package instruqt

import (
	"context"
	"time"

	"github.com/shurcooL/graphql"
)

type inviteQuery struct {
	TrackInvite `graphql:"trackInvite(inviteID: $inviteId)"`
}

type TrackInvite struct {
	Id                string
	PublicTitle       string
	RuntimeParameters struct {
		EnvironmentVariables []variable
	}
	Claims []TrackInviteClaim
}

type TrackInviteClaim struct {
	Id   string
	User struct {
		Id string
	}
	ClaimedAt time.Time
}

type variable struct {
	Key   string
	Value string
}

func (c *Client) GetInvite(inviteId string) (i TrackInvite, err error) {
	if inviteId == "" {
		return i, nil
	}

	var q inviteQuery
	variables := map[string]interface{}{
		"inviteId": graphql.String(inviteId),
	}

	ctx := context.Background()
	if err := c.GraphQLClient.Query(ctx, &q, variables); err != nil {
		return i, err
	}

	return q.TrackInvite, nil
}

type invitesQuery struct {
	TrackInvites []TrackInvite `graphql:"trackInvites(teamSlug: $teamSlug)"`
}

func (c *Client) GetInvites() (i []TrackInvite, err error) {
	var q invitesQuery
	variables := map[string]interface{}{
		"teamSlug": graphql.String(teamSlug),
	}

	ctx := context.Background()
	if err := c.GraphQLClient.Query(ctx, &q, variables); err != nil {
		return i, err
	}

	return q.TrackInvites, nil
}
