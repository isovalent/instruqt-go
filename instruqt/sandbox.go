package instruqt

import (
	"time"

	"github.com/shurcooL/graphql"
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
	Last_Activity_At time.Time    // The timestamp of the last activity in the sandbox.
	State            string       // The current state of the sandbox (e.g., "running", "stopped").
	Track            SandboxTrack // The track associated with the sandbox.
	Invite           TrackInvite  // The invite details associated with the sandbox.
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
