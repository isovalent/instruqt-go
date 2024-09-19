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

	"github.com/shurcooL/graphql"
)

// teamQuery represents the GraphQL query structure for retrieving the TPG public key
// associated with a specific team identified by its slug.
type teamQuery struct {
	Team struct {
		TPGPublicKey graphql.String `graphql:"tpgPublicKey"` // The TPG public key of the team.
	} `graphql:"team(teamSlug: $teamSlug)"`
}

// GetTPGPublicKey retrieves the TPG public key for the team associated with the client.
//
// Returns:
//   - string: The TPG public key of the team.
//   - error: Any error encountered while retrieving the TPG public key.
func (c *Client) GetTPGPublicKey() (string, error) {
	var q teamQuery
	variables := map[string]interface{}{
		"teamSlug": graphql.String(c.TeamSlug),
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return "", fmt.Errorf("failed to retrieve TPG Public Key: %v", err)
	}

	return string(q.Team.TPGPublicKey), nil
}
