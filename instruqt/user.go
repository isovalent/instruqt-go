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
	"strings"

	graphql "github.com/hasura/go-graphql-client"
)

// userInfoQuery represents the GraphQL query structure for retrieving user information
// by the user's unique ID.
type userInfoQuery struct {
	User `graphql:"user(userID: $userID)"`
}

// User represents the data structure for an Instruqt user.
type User struct {
	Id      string
	Details *UserDetails `graphql:"details(teamSlug: $teamSlug)"`
	Profile *UserProfile
}

// Detailed user information associated with a specific team.
type UserDetails struct {
	FirstName   graphql.String  // The first name of the user.
	LastName    graphql.String  // The last name of the user.
	Email       graphql.String  // The email of the user.
	CompanyName graphql.String  // The company name of the user.
	JobTitle    graphql.String  // The job title of the user.
	JobLevel    graphql.String  // The job level of the user.
	CountryCode graphql.String  // The country code of the user.
	Consent     graphql.Boolean // The consent status of the user.
}

// Profile-level information for the user.
type UserProfile struct {
	Display_Name graphql.String // The display name of the user.
	Email        graphql.String // The email of the user.
}

// UserInfo represents a simplified user information structure.
type UserInfo struct {
	FirstName string // The first name of the user.
	LastName  string // The last name of the user.
	Email     string // The email of the user.
}

// GetUserInfo retrieves the user information from Instruqt using the user's unique ID.
//
// Parameters:
//   - userId: The unique identifier of the user.
//
// Returns:
//   - UserInfo: The user's information including first name, last name, and email.
//   - error: Any error encountered while retrieving the user information.
func (c *Client) GetUserInfo(userId string) (u UserInfo, err error) {
	var q userInfoQuery
	variables := map[string]interface{}{
		"teamSlug": graphql.String(c.TeamSlug),
		"userID":   graphql.String(userId),
	}
	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return u, fmt.Errorf("[GetUserInfo] Failed to retrieve user info: %v", err)
	}

	if q.User.Details != nil && q.User.Details.Email != "" {
		c.InfoLogger.Printf("[Instruqt][GetUserInfo][%s] Found user info from instruqt user details", userId)
		u = UserInfo{
			FirstName: string(q.User.Details.FirstName),
			LastName:  string(q.User.Details.LastName),
			Email:     string(q.User.Details.Email),
		}
		return u, nil
	}

	if q.User.Profile != nil && q.User.Profile.Email != "" {
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
