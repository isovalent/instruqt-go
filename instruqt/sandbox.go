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
	"time"

	graphql "github.com/hasura/go-graphql-client"
)

// sandboxVarQuery represents the GraphQL query structure for retrieving a single
// sandbox variable based on the sandbox ID, hostname, and key.
type sandboxVarQuery struct {
	GetSandboxVariable SandboxVar `graphql:"getSandboxVariable(sandboxID: $sandboxID, hostname: $hostname, key: $key)"`
}

// SandboxVar represents a key-value pair for a variable within a sandbox environment.
type SandboxVar struct {
	Key   string // The key of the sandbox variable.
	Value string // The value of the sandbox variable.
}

// sandboxQuery represents the GraphQL query structure for a single sandbox by its ID
type sandboxQuery struct {
	Sandbox Sandbox `graphql:"sandbox(ID: $id)"`
}

// sandboxesQuery represents the GraphQL query structure for retrieving all sandboxes
// associated with a specific team.
type sandboxesQuery struct {
	Sandboxes struct {
		Nodes []Sandbox // A list of sandboxes retrieved by the query.
	} `graphql:"sandboxes(teamSlug: $teamSlug)"`
}

// Sandbox represents a sandbox environment within Instruqt, including details
// about its state, associated track, and invite.
type Sandbox struct {
	Id               string        // The id of the sandbox.
	Last_Activity_At time.Time     // The timestamp of the last activity in the sandbox.
	State            string        // The current state of the sandbox (e.g., "running", "stopped").
	Track            SandboxTrack  // The track associated with the sandbox.
	Invite           TrackInvite   // The invite details associated with the sandbox.
	User             User          // The user running the sandbox.
	Hot_Start_Pool   *HotStartPool // The hot start pool associated with the sandbox.
}

// GetSandboxVariable retrieves a specific variable from a sandbox environment
// using the sandbox ID and the variable's key.
//
// Parameters:
//   - playID: The unique identifier of the sandbox environment.
//   - key: The key of the sandbox variable to retrieve.
//
// Returns:
//   - string: The value of the requested sandbox variable.
//   - error: Any error encountered while retrieving the variable.
func (c *Client) GetSandboxVariable(playID string, key string) (v string, err error) {
	if playID == "" || key == "" {
		return v, nil
	}

	var hostname = "server"

	var q sandboxVarQuery
	variables := map[string]interface{}{
		"hostname":  graphql.String(hostname),
		"sandboxID": graphql.String(playID),
		"key":       graphql.String(key),
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return v, err
	}

	return q.GetSandboxVariable.Value, nil
}

// GetSandbox retrieves a sandbox by its ID.
//
// Returns:
//   - Sandbox: The sandbox.
//   - error: Any error encountered while retrieving the sandbox.
func (c *Client) GetSandbox(id string) (s Sandbox, err error) {
	var q sandboxQuery
	variables := map[string]interface{}{
		"id":       graphql.ID(id),
		"teamSlug": graphql.String(c.TeamSlug), // Pass teamSlug for User info
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return s, err
	}

	return q.Sandbox, nil
}

// GetSandboxes retrieves all sandboxes associated with the team slug defined in the client.
//
// Returns:
//   - []Sandbox: A list of sandboxes for the team.
//   - error: Any error encountered while retrieving the sandboxes.
func (c *Client) GetSandboxes() (s []Sandbox, err error) {
	var q sandboxesQuery
	variables := map[string]interface{}{
		"teamSlug": graphql.String(c.TeamSlug),
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return s, err
	}

	return q.Sandboxes.Nodes, nil
}
