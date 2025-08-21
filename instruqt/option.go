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

// instruqt/options.go

package instruqt

// Option defines a functional option for configuring methods.
// It allows modifying the behavior of query methods, such as including additional fields.
type Option func(*options)

// options holds configuration options for all methods.
// It includes flags and parameters that can be used by different methods.
type options struct {
	// Options for GetReview
	includePlay bool

	// Options for GetTrack*
	includeChallenges bool
	includeReviews    bool

	// Options for GetChallenge*
	includeAssignment bool

	// Options for GetPlays
	trackIDs       []string
	trackInviteIDs []string
	landingPageIDs []string
	tags           []string
	userIDs        []string
	playType       PlayType
	ordering       *Ordering

	// Options for GetSandboxes
	state   string
	poolIDs []string
}

// WithPlay is a functional option that configures methods to include the 'play' field in the query.
// Example usage: GetReview("reviewID", WithPlay())
func WithPlay() Option {
	return func(opts *options) {
		opts.includePlay = true
	}
}

// WithChallenges is a functional option to include challenges.
// Example usage: GetTrackById("tracKID", WithChallenges())
func WithChallenges() Option {
	return func(opts *options) {
		opts.includeChallenges = true
	}
}

// WithReviews is a functional option to include reviews.
// Example usage: GetTrackById("tracKID", WithReviews())
func WithReviews() Option {
	return func(opts *options) {
		opts.includeReviews = true
	}
}

// WithTrackIDs sets the TrackIDs filter for methods that support it.
// Usage: GetPlays(from, to, take, skip, WithTrackIDs("track1", "track2"))
func WithTrackIDs(ids ...string) Option {
	return func(opts *options) {
		opts.trackIDs = ids
	}
}

// WithTrackInviteIDs sets the TrackInviteIDs filter for methods that support it.
// Usage: GetPlays(from, to, take, skip, WithTrackInviteIDs("invite1", "invite2"))
func WithTrackInviteIDs(ids ...string) Option {
	return func(opts *options) {
		opts.trackInviteIDs = ids
	}
}

// WithTags sets the Tags filter for methods that support it.
// Usage: GetPlays(from, to, take, skip, WithTags("tag1", "tag2"))
func WithTags(tags ...string) Option {
	return func(opts *options) {
		opts.tags = tags
	}
}

// WithUserIDs sets the UserIDs filter for methods that support it.
// Usage: GetPlays(from, to, take, skip, WithUserIDs("user1", "user2"))
func WithUserIDs(ids ...string) Option {
	return func(opts *options) {
		opts.userIDs = ids
	}
}

// WithPlayType sets the PlayType filter for methods that support it.
// Usage: GetPlays(from, to, take, skip, WithPlayType(PlayTypeDeveloper))
func WithPlayType(pt PlayType) Option {
	return func(opts *options) {
		opts.playType = pt
	}
}

// WithState sets the State filter for methods that support it.
// Usage: GetSandboxes(WithState("active"))
func WithState(state string) Option {
	return func(opts *options) {
		opts.state = state
	}
}

// WithPoolIDs sets the PoolIDs filter for methods that support it.
// Usage: GetSandboxes(WithPoolIDs("pool1", "pool2"))
func WithPoolIDs(ids ...string) Option {
	return func(opts *options) {
		opts.poolIDs = ids
	}
}

func WithAssignment() Option {
	return func(opts *options) {
		opts.includeAssignment = true
	}
}

// OrderBy represents the fields by which plays can be ordered.
type OrderBy string

const (
	OrderByCompletionPercent OrderBy = "completion_percent"
	OrderByTimeSpent         OrderBy = "time_spent"
)

// Direction represents the sorting direction.
type Direction string

const (
	DirectionAsc  Direction = "Asc"
	DirectionDesc Direction = "Desc"
)

// Ordering represents the sorting parameters for plays.
type Ordering struct {
	OrderBy   OrderBy   // Must be "completion_percent" or "time_spent"
	Direction Direction // "Asc" or "Desc"
}

// WithOrdering sets the ordering parameters for methods that support it.
// Usage: GetPlays(from, to, take, skip, WithOrdering(OrderByCompletionPercent, DirectionDesc))
func WithOrdering(orderBy OrderBy, direction Direction) Option {
	return func(opts *options) {
		opts.ordering = &Ordering{
			OrderBy:   orderBy,
			Direction: direction,
		}
	}
}
