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

// sandboxVarSet represents the GraphQL mutation structure for setting a single
type sandboxVarSet struct {
	SetSandboxVariable SandboxVar `graphql:"setSandboxVariable(sandboxID: $sandboxID, hostname: $hostname, key: $key, value: $value)"`
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

// SandboxState defines the possible states of a sandbox.
type SandboxState string

// Constants representing different sandbox states.
const (
	SandboxStateCreating SandboxState = "creating"
	SandboxStateCreated  SandboxState = "created"
	SandboxStateFailed   SandboxState = "failed"
	SandboxStatePooled   SandboxState = "pooled"
	SandboxStateStopped  SandboxState = "stopped"
	SandboxStateActive   SandboxState = "active"
	SandboxStateClaimed  SandboxState = "claimed"
	SandboxStateCleaning SandboxState = "cleaning"
	SandboxStateCleaned  SandboxState = "cleaned"
)

// associated with a specific team.
type sandboxesQuery struct {
	Sandboxes struct {
		Nodes []Sandbox // A list of sandboxes retrieved by the query.
	} `graphql:"sandboxes(teamSlug: $teamSlug, filter: {track_ids: $track_ids, invite_ids: $invite_ids, pool_ids: $pool_ids, user_name_or_id: $user_name_or_id, state: $state})"`
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
func (c *Client) GetSandboxVariable(playID string, hostname string, key string) (v string, err error) {
	if playID == "" || key == "" {
		return v, nil
	}

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

// SetSandboxVariable sets a specific variable in a sandbox environment
// using the sandbox ID, variable key, and value.
func (c *Client) SetSandboxVariable(playID string, hostname string, key string, value string) error {
	if playID == "" || key == "" || value == "" {
		return nil
	}

	var q sandboxVarSet
	variables := map[string]interface{}{
		"hostname":  graphql.String(hostname),
		"sandboxID": graphql.String(playID),
		"key":       graphql.String(key),
		"value":     graphql.String(value),
	}

	if err := c.GraphQLClient.Mutate(c.Context, &q, variables); err != nil {
		return err
	}

	return nil
}

// GetSandbox retrieves a sandbox by its ID.
//
// Returns:
//   - Sandbox: The sandbox.
//   - error: Any error encountered while retrieving the sandbox.
func (c *Client) GetSandbox(id string, opts ...Option) (s Sandbox, err error) {
	// Initialize the filter with default values
	filters := &options{
		playType: PlayTypeAll, // Default PlayType
	}

	// Apply each option to modify the filter
	for _, opt := range opts {
		opt(filters)
	}

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
func (c *Client) GetSandboxes(opts ...Option) (s []Sandbox, err error) {
	// Initialize the filter with default values
	filters := &options{
		playType: PlayTypeAll, // Default PlayType
	}

	// Apply each option to modify the filter
	for _, opt := range opts {
		opt(filters)
	}

	// Convert Go types to GraphQL types
	trackIds := make([]graphql.String, len(filters.trackIDs))
	for i, id := range filters.trackIDs {
		trackIds[i] = graphql.String(id)
	}

	trackInviteIds := make([]graphql.String, len(filters.trackInviteIDs))
	for i, id := range filters.trackInviteIDs {
		trackInviteIds[i] = graphql.String(id)
	}

	poolIds := make([]graphql.String, len(filters.poolIDs))
	for i, id := range filters.poolIDs {
		poolIds[i] = graphql.String(id)
	}

	var userNameOrId string
	if len(filters.userIDs) > 0 {
		userNameOrId = filters.userIDs[0]
	}

	var q sandboxesQuery
	variables := map[string]interface{}{
		"teamSlug":        graphql.String(c.TeamSlug),
		"track_ids":       trackIds,
		"invite_ids":      trackInviteIds,
		"pool_ids":        poolIds,
		"user_name_or_id": graphql.String(userNameOrId),
		"state":           filters.states,
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return s, err
	}

	return q.Sandboxes.Nodes, nil
}
