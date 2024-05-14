package instruqt

import (
	"context"

	"github.com/shurcooL/graphql"
)

type sandboxVarQuery struct {
	GetSandboxVariable SandboxVar `graphql:"getSandboxVariable(sandboxID: $sandboxID, hostname: $hostname, key: $key)"`
}

type SandboxVar struct {
	Key   string
	Value string
}

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

	ctx := context.Background()
	if err := c.GraphQLClient.Query(ctx, &q, variables); err != nil {
		return v, err
	}

	return q.GetSandboxVariable.Value, nil
}
