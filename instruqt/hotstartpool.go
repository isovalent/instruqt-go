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
	"time"

	graphql "github.com/hasura/go-graphql-client"
)

// HotStartPoolType defines a custom type for HotStartPool types.
type HotStartPoolType string

// Constants representing different types of HotStartPool.
const (
	HotStartPoolTypeDedicated HotStartPoolType = "dedicated"
	HotStartPoolTypeShared    HotStartPoolType = "shared"
)

// HotStartStatus defines a custom type for HotStartPool status.
type HotStartStatus string

// Constats representing the different types of status of HotStartPool.
const (
	HostStartStatusRunning      HotStartStatus = "Running"
	HostStartStatusProvisioning HotStartStatus = "Provisioning"
	HostStartStatusInactive     HotStartStatus = "Inactive"
	HostStartStatusExpired      HotStartStatus = "Expired"
	HostStartStatusDeleted      HotStartStatus = "Deleted"
	HostStartStatusAutoRefill   HotStartStatus = "AutoRefill"
)

// HotStartPoolConfigTrackEdge
type HotStartPoolConfigTrackEdge struct {
	Claimed   int
	Available int
	Created   int
	Failed    int
	Creating  int
	Total     int
	Node      SandboxConfig
}

// HotStartPoolTrackEdge
type HotStartPoolTrackEdge struct {
	Claimed   int
	Available int
	Created   int
	Failed    int
	Creating  int
	Total     int
	Node      Track
}

// HotStartPool represents a hot start pool in Instruqt.
type HotStartPool struct {
	Id          string                        // ID of the hot start pool.
	Type        HotStartPoolType              // The type of hot start pool.
	Size        int                           // Number of sandboxes available per track.
	Created     *time.Time                    // Creation time of the hot start pool.
	Deleted     *time.Time                    // Deletion time of the hot start pool.
	Name        string                        // Name given to the hot start pool.
	Auto_refill bool                          // Flag that signals if the sandboxes should be auto refillable.
	Starts_at   *time.Time                    // Schedule time for the hot start pool to start creating sandboxes.
	Ends_at     *time.Time                    // Schedule time for the hot start pool to stop creating sandboxes.
	Status      HotStartStatus                // Status of the hot start pool.
	Region      string                        // Region of a hotstart pool.
	Configs     []HotStartPoolConfigTrackEdge // Configs status for the hotstart pool.
	Tracks      []HotStartPoolTrackEdge       // Tracks status for the hotstart pool.
}

type hotStartPoolConnection struct {
	Nodes      []hotStartPoolQueryNode
	TotalCount int
	PageInfo   struct {
		EndCursor   string
		HasNextPage bool
	}
}

type hotStartPoolQueryNode struct {
	Id          string
	Type        HotStartPoolType
	Size        int
	Created     *time.Time
	Deleted     *time.Time
	Name        string
	Auto_refill bool
	Starts_at   *time.Time
	Ends_at     *time.Time
	Status      HotStartStatus
	Region      string
	Configs     []struct {
		Claimed   int
		Available int
		Created   int
		Failed    int
		Creating  int
		Total     int
		Node      struct {
			Id      string
			Name    string
			Slug    string
			Version int
		}
	}
	Tracks []struct {
		Claimed   int
		Available int
		Created   int
		Failed    int
		Creating  int
		Total     int
		Node      struct {
			Id    string
			Slug  string
			Title string
		}
	}
}

type hotStartPoolsFirstPageQuery struct {
	HotStartPools hotStartPoolConnection `graphql:"hotStartPools(organizationSlug: $organizationSlug, teamSlug: $teamSlug, paging: {First: $first})"`
}

type hotStartPoolsAfterPageQuery struct {
	HotStartPools hotStartPoolConnection `graphql:"hotStartPools(organizationSlug: $organizationSlug, teamSlug: $teamSlug, paging: {First: $first, After: $after})"`
}

// GetHotStartPools retrieves a page of hot start pools for the client's
// organization and team. The returned cursor is used to request subsequent
// pages from the connection.
func (c *Client) GetHotStartPools(first int, after string) ([]HotStartPool, int, string, bool, error) {
	if first <= 0 {
		first = 100
	}

	variables := map[string]interface{}{
		"organizationSlug": graphql.String(c.TeamSlug),
		"teamSlug":         graphql.String(c.TeamSlug),
		"first":            graphql.Int(first),
	}
	if after == "" {
		var q hotStartPoolsFirstPageQuery
		if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
			return nil, 0, "", false, fmt.Errorf("GraphQL query failed: %w", err)
		}
		return convertHotStartPools(q.HotStartPools.Nodes), q.HotStartPools.TotalCount, q.HotStartPools.PageInfo.EndCursor, q.HotStartPools.PageInfo.HasNextPage, nil
	}

	var q hotStartPoolsAfterPageQuery
	variables["after"] = graphql.String(after)
	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return nil, 0, "", false, fmt.Errorf("GraphQL query failed: %w", err)
	}
	return convertHotStartPools(q.HotStartPools.Nodes), q.HotStartPools.TotalCount, q.HotStartPools.PageInfo.EndCursor, q.HotStartPools.PageInfo.HasNextPage, nil
}

func convertHotStartPools(nodes []hotStartPoolQueryNode) []HotStartPool {
	pools := make([]HotStartPool, 0, len(nodes))
	for _, node := range nodes {
		pool := HotStartPool{
			Id: node.Id, Type: node.Type, Size: node.Size, Created: node.Created, Deleted: node.Deleted,
			Name: node.Name, Auto_refill: node.Auto_refill, Starts_at: node.Starts_at, Ends_at: node.Ends_at,
			Status: node.Status, Region: node.Region,
		}
		for _, config := range node.Configs {
			pool.Configs = append(pool.Configs, HotStartPoolConfigTrackEdge{
				Claimed: config.Claimed, Available: config.Available, Created: config.Created, Failed: config.Failed,
				Creating: config.Creating, Total: config.Total,
				Node: SandboxConfig{Id: config.Node.Id, Name: config.Node.Name, Slug: config.Node.Slug, Version: config.Node.Version},
			})
		}
		for _, track := range node.Tracks {
			pool.Tracks = append(pool.Tracks, HotStartPoolTrackEdge{
				Claimed: track.Claimed, Available: track.Available, Created: track.Created, Failed: track.Failed,
				Creating: track.Creating, Total: track.Total,
				Node: Track{Id: track.Node.Id, Slug: track.Node.Slug, Title: track.Node.Title},
			})
		}
		pools = append(pools, pool)
	}
	return pools
}
