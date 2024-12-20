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

type SandboxConfig struct {
	Id      string
	Name    string
	Slug    string
	Version int
	Deleted *time.Time
}

type SandboxConfigVersionStatus string

const (
	SandboxConfigVersionStatusDraft     SandboxConfigVersionStatus = "draft"
	SandboxConfigVersionStatusPublished SandboxConfigVersionStatus = "published"
	SandboxConfigVersionStatusArchived  SandboxConfigVersionStatus = "archived"
	SandboxConfigVersionStatusIsolated  SandboxConfigVersionStatus = "isolated"
)

type SandboxConfigVersion struct {
	Id           string
	Config       SandboxConfig
	Version      int
	Description  string
	Status       SandboxConfigVersionStatus
	Published_At *time.Time
}
