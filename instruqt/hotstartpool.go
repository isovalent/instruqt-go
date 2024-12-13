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

import "time"

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
