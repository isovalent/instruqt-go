package instruqt

import (
	"context"
	"fmt"
	"strings"

	"github.com/shurcooL/graphql"
)

type userInfoQuery struct {
	User `graphql:"user(userID: $userID)"`
}

type User struct {
	Details struct {
		FirstName graphql.String
		LastName  graphql.String
		Email     graphql.String
	} `graphql:"details(teamSlug: $teamSlug)"`
	Profile struct {
		Display_Name graphql.String
		Email        graphql.String
	}
}

type UserInfo struct {
	FirstName string
	LastName  string
	Email     string
}

func (c *Client) GetUserInfo(userId string) (u UserInfo, err error) {
	var q userInfoQuery
	ctx := context.Background()
	variables := map[string]interface{}{
		"teamSlug": graphql.String("isovalent"),
		"userID":   graphql.String(userId),
	}
	if err := c.GraphQLClient.Query(ctx, &q, variables); err != nil {
		return u, fmt.Errorf("[GetUserInfo] Failed to retrieve user info: %v", err)
	}

	if q.User.Details.Email != "" {
		c.InfoLogger.Printf("[Instruqt][GetUserInfo][%s] Found user info from instruqt user details", userId)
		u = UserInfo{
			FirstName: string(q.User.Details.FirstName),
			LastName:  string(q.User.Details.LastName),
			Email:     string(q.User.Details.Email),
		}
		return u, nil
	}

	if q.User.Profile.Email != "" {
		c.InfoLogger.Printf("[Instruqt][GetUserInfo][%s] Found user info from instruqt user profile", userId)
		nameParts := strings.Fields(string(q.User.Profile.Display_Name))
		u = UserInfo{
			FirstName: nameParts[0],
			LastName:  strings.Join(nameParts[1:], " "),
			Email:     string(q.User.Profile.Email),
		}

		return u, nil
	}

	return u, nil
}
